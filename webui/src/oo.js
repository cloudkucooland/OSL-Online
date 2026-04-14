export const server = 'https://saint-luke.net:8443';

/**
 * Handles the heavy lifting for all API requests to ensure consistency.
 */
async function apiCall(endpoint, options = {}) {
	const jwt = localStorage.getItem('jwt');

	const defaultOptions = {
		mode: 'cors',
		credentials: 'include',
		referrerPolicy: 'origin',
		headers: {}
	};

	// Merge headers
	if (jwt) {
		defaultOptions.headers['Authorization'] = `Bearer ${jwt}`;
	}

	const finalOptions = { ...defaultOptions, ...options };

	const response = await fetch(`${server}/api/v1${endpoint}`, finalOptions);

	if (response.status === 401) {
		localStorage.removeItem('jwt');
		window.location.hash = '#/login';
		throw new Error('Session expired');
	}

	if (!response.ok) {
		const errorData = await response.json().catch(() => ({ error: 'Unknown server error' }));
		console.error(`Server returned ${response.status}:`, errorData);
		throw new Error(errorData.error || 'Server error');
	}

	return response;
}

export function cleanDateFormat(incoming) {
	if (!incoming) return '';
	const d = new Date(incoming);
	if (isNaN(d.valueOf())) return incoming;

	const fmt = d.toISOString().slice(0, 10);
	return fmt === '0001-01-01' ? '' : fmt;
}

const MEMBER_DATE_FIELDS = [
	'DateReaffirmation',
	'DateFirstVows',
	'DateNovitiate',
	'BirthDate',
	'DateDeceased',
	'DateRecordCreated',
	'DateRemoved',
	'DateLifeVows'
];

function processMemberDates(member) {
	MEMBER_DATE_FIELDS.forEach((field) => {
		if (member[field]) member[field] = cleanDateFormat(member[field]);
	});
	return member;
}

export function oslname(m) {
	let name = '';
	const isVowed = m.MemberStatus === 'Life Vows' || m.MemberStatus === 'Annual Vows';

	if (isVowed) name = `${m.Title} `;

	let firstName = m.PreferredName || m.FirstName;

	if (m.MemberStatus === 'Life Vows' && m.LifevowName) {
		name += `${m.LifevowName} `;
	} else {
		name += `${firstName} `;
	}

	name += m.LastName;
	if (isVowed) name += ', OSL';

	return name;
}

export function getMe() {
	const jwt = localStorage.getItem('jwt');
	if (!jwt) undefined;

	try {
		const token = JSON.parse(window.atob(jwt.split('.')[1]));
		const exp = new Date(token.exp * 1000);

		if (exp <= Date.now()) {
			localStorage.removeItem('jwt');
			return undefined;
		}

		// Refresh if expiring in the next week
		if (exp.valueOf() - 604800000 <= Date.now()) {
			refreshJWT();
		}
		return token;
	} catch (e) {
		return undefined;
	}
}

export async function getJWT(username, password) {
	const body = new FormData();
	body.append('username', username);
	body.append('password', password);

	const res = await apiCall('/getJWT', { method: 'POST', body });
	const payload = await res.text();
	localStorage.setItem('jwt', payload);
	try {
		const token = JSON.parse(window.atob(jwt.split('.')[1]));
		const exp = new Date(token.exp * 1000);

		if (exp <= Date.now()) {
			console.log('server sent expired token!');
			// return undefined
		}
		if (token.aud[0] != 'OSL-Online') {
			console.log('server sent weird token!');
			// return undefined
		}
	} catch (e) {
		console.log('unable to parse token');
		return undefined;
	}
	return payload;
}

export async function refreshJWT() {
	const res = await apiCall('/refreshJWT', { method: 'GET' });
	const payload = await res.text();
	localStorage.setItem('jwt', payload);
}

export async function search(query) {
	const body = new FormData();
	body.append('query', query);
	const res = await apiCall('/search', { method: 'POST', body });
	return await res.json();
}

export async function getMember(id) {
	const res = await apiCall(`/member/${id}`);
	const payload = await res.json();
	return processMemberDates(payload);
}

export async function updateMember(id, field, value) {
	const body = new FormData();
	body.append('field', field);
	body.append('value', value);
	await apiCall(`/member/${id}`, { method: 'PUT', body });
}

export async function createMember(firstname, lastname) {
	const body = new FormData();
	body.append('firstname', firstname);
	body.append('lastname', lastname);
	const res = await apiCall('/member', { method: 'POST', body });
	const payload = await res.json();
	return payload.id;
}

// --- Subscriber ---

export async function subsearch(query) {
	const body = new FormData();
	body.append('query', query);
	const res = await apiCall('/subsearch', { method: 'POST', body });
	return await res.json();
}

export async function getSubscriber(id) {
	const res = await apiCall(`/subscriber/${id}`);
	const payload = await res.json();
	if (payload.DatePaid) payload.DatePaid = cleanDateFormat(payload.DatePaid);
	return payload;
}

export async function updateSubscriber(id, field, value) {
	const body = new FormData();
	body.append('field', field);
	body.append('value', value);
	await apiCall(`/subscriber/${id}`, { method: 'PUT', body });
}

// --- Reports & Exports ---

export async function report(reportname) {
	const res = await apiCall(`/report/${reportname}`);
	const blob = await res.blob();
	const contentType = res.headers.get('Content-Type') || '';
	const extension = contentType.toLowerCase().includes('pdf') ? 'pdf' : 'csv';

	const url = URL.createObjectURL(blob);
	const link = document.createElement('a');
	link.href = url;
	link.download = `${reportname}.${extension}`;
	document.body.appendChild(link);
	link.click();
	document.body.removeChild(link);
	URL.revokeObjectURL(url);
}

export async function vcard(id) {
	const res = await apiCall(`/member/${id}/vcard`);
	const blob = await res.blob();
	const url = URL.createObjectURL(blob);
	const link = document.createElement('a');
	link.href = url;
	link.download = `${id}.vcf`;
	document.body.appendChild(link);
	link.click();
	document.body.removeChild(link);
	URL.revokeObjectURL(url);
}

// --- Auth User Profile ---

export async function getMeFromServer() {
	const res = await apiCall('/me');
	const payload = await res.json();
	return processMemberDates(payload);
}

export async function updateMe(field, value) {
	const body = new FormData();
	body.append('field', field);
	body.append('value', value);
	await apiCall('/me', { method: 'PUT', body });
}

export async function postRegister(email) {
	const body = new FormData();
	body.append('email', email);
	await apiCall('/register', { method: 'POST', body });
	return true;
}

// --- Financials / Giving ---

export async function getMeGiving() {
	const res = await apiCall('/me/giving');
	return await res.json();
}

export async function getGiving(id) {
	const res = await apiCall(`/giving/${id}`);
	const payload = await res.json();
	return payload.map((row) => ({ ...row, Date: cleanDateFormat(row.Date) }));
}

export async function postGiving(id, date, amount, description, check, transaction) {
	const body = new FormData();
	body.append('id', id);
	body.append('date', cleanDateFormat(date));
	body.append('amount', amount);
	body.append('description', description);
	body.append('check', check);
	body.append('transaction', transaction);
	await apiCall(`/giving/${id}`, { method: 'POST', body });
}

// --- Logs & Chapters ---

export async function getChangelog(id) {
	const res = await apiCall(`/changelog/${id}`);
	const payload = await res.json();
	return payload.map((row) => ({ ...row, Date: cleanDateFormat(row.Date) }));
}

export async function getChapters() {
	console.log('getChapters');
	const res = await apiCall('/chapter');
	const payload = await res.json();
	return payload.map((c) => ({ ...c, value: c.ID, name: c.Name }));
}

export async function getMeChapters() {
	const res = await apiCall('/me/chapters');
	return await res.json();
}

export async function updateMeChapters(chapters) {
	const body = new FormData();
	body.append('chapters', chapters);
	const res = await apiCall('/me/chapters', { method: 'PUT', body });
	return await res.json();
}

export async function getMemberChapters(id) {
	const res = await apiCall(`/member/${id}/chapters`);
	return await res.json();
}

export async function updateMemberChapters(id, chapters) {
	const body = new FormData();
	body.append('chapters', chapters);
	const res = await apiCall(`/member/${id}/chapters`, { method: 'PUT', body });
	return await res.json();
}

export async function getChapterMembers(chapterID) {
	const res = await apiCall(`/chapter/${chapterID}`);
	return await res.json();
}

// --- Admin / Misc ---

export async function getLeaders(category) {
	const res = await apiCall(`/leaders/${category}`);
	return await res.json();
}

export async function searchemail(query) {
	const body = new FormData();
	body.append('query', query);
	const res = await apiCall('/searchemail', { method: 'POST', body });
	return await res.json();
}

export async function getLocalities() {
	const res = await apiCall('/localities');
	const payload = await res.json();
	return payload.map((c) => ({
		...c,
		value: c.JointCode,
		name: `${c.CountryCode}: ${c.Locality}`
	}));
}

export async function getLocalityMembers(loc) {
	const res = await apiCall(`/locality/${loc}`);
	return await res.json();
}

export async function sendemail(whom, subject, message) {
	const body = new FormData();
	body.append('message', message);
	body.append('whom', whom);
	body.append('subject', subject);
	await apiCall('/email', { method: 'POST', body });
	return true;
}

export async function getDashboard() {
	const res = await apiCall('/dashboard');
	return await res.json();
}

export async function getNecrology() {
	const res = await apiCall('/necrology');
	return await res.json();
}

// --- Notes ---

export async function getMemberNotes(id) {
	const res = await apiCall(`/member/${id}/notes`);
	return await res.json();
}

export async function postMemberNote(id, note) {
	const body = new FormData();
	body.append('note', note);
	await apiCall(`/member/${id}/notes`, { method: 'POST', body });
	return true;
}

export async function deleteMemberNote(memberid, noteid) {
	await apiCall(`/member/${memberid}/notes/${noteid}`, { method: 'DELETE' });
	return true;
}

export async function getAllPrayers() {
	const res = await apiCall(`/prayers`);
	return await res.json();
}

export async function getMyPrayers(memberid) {
	const res = await apiCall(`/member/${memberid}/prayers`);
	return await res.json();
}

export async function addPrayer(prayerText, isAnonymous = false) {
	const res = await apiCall(`/prayers`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ Content: prayerText, Anonymous: isAnonymous })
	});
	return res.ok;
}

export async function deletePrayer(prayerid) {
	console.log('deleting prayer', prayerid);
	const res = await apiCall(`/prayers/${prayerid}`, {
		method: 'DELETE'
	});
	return res.ok;
}
