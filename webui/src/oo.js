export const server = 'https://saint-luke.net:8443';

export function getMe() {
	const jwt = localStorage.getItem('jwt');
	if (jwt === undefined || jwt === null) {
		return undefined;
	}

	const token = JSON.parse(window.atob(jwt.split('.')[1]).toString());
	return token;
}

export async function getJWT(username, password) {
	const dataArray = new FormData();
	dataArray.append('username', username);
	dataArray.append('password', password);

	const request = {
		method: 'POST',
		mode: 'cors',
		credentials: 'include',
		redirect: 'manual',
		referrerPolicy: 'origin',
		body: dataArray
	};

	const response = await fetch(`${server}/api/v1/getJWT`, request);
	const payload = await response.text();
	if (response.status != 200) {
		console.log('server returned ', response.status);
		throw new Error(payload.error);
	}
	localStorage.setItem('jwt', payload);
}

export async function search(query) {
	const jwt = localStorage.getItem('jwt');
	if (jwt === undefined || jwt === null) {
		throw new Error('Not Logged in');
	}

	const dataArray = new FormData();
	dataArray.append('query', query);

	const request = {
		method: 'POST',
		mode: 'cors',
		credentials: 'include',
		redirect: 'manual',
		referrerPolicy: 'origin',
		body: dataArray,
		headers: { Authorization: 'Bearer ' + jwt }
	};

	const response = await fetch(`${server}/api/v1/search`, request);
	const payload = await response.json();
	if (response.status != 200) {
		console.log('server returned ', response.status);
		throw new Error(payload.error);
	}
	// console.log(payload);
	return payload;
}

export async function getMember(id) {
	const jwt = localStorage.getItem('jwt');
	if (jwt === undefined || jwt === null) {
		throw new Error('Not Logged in');
	}

	const request = {
		method: 'GET',
		mode: 'cors',
		credentials: 'include',
		redirect: 'manual',
		referrerPolicy: 'origin',
		headers: { Authorization: 'Bearer ' + jwt }
	};

	const response = await fetch(`${server}/api/v1/member/${id}`, request);
	const payload = await response.json();
	if (response.status != 200) {
		console.log('server returned ', response.status);
		throw new Error(payload.error);
	}

	// trim unnecessary times
	let sp = payload.DateReaffirmation.split('T');
	if (sp[0] == '0001-01-01') {
		payload.DateReaffirmation = '';
	} else {
		payload.DateReaffirmation = sp[0];
	}

	sp = payload.DateFirstVows.split('T');
	if (sp[0] == '0001-01-01') {
		payload.DateFirstVows = '';
	} else {
		payload.DateFirstVows = sp[0];
	}

	sp = payload.DateNovitiate.split('T');
	if (sp[0] == '0001-01-01') {
		payload.DateNovitiate = '';
	} else {
		payload.DateNovitiate = sp[0];
	}

	sp = payload.BirthDate.split('T');
	if (sp[0] == '0001-01-01') {
		payload.BirthDate = '';
	} else {
		payload.BirthDate = sp[0];
	}

	sp = payload.DateDeceased.split('T');
	if (sp[0] == '0001-01-01') {
		payload.DateDeceased = '';
	} else {
		payload.DateDeceased = sp[0];
	}

	sp = payload.DateRecordCreated.split('T');
	if (sp[0] == '0001-01-01') {
		payload.DateRecordCreated = '';
	} else {
		payload.DateRecordCreated = sp[0];
	}

	sp = payload.DateRemoved.split('T');
	if (sp[0] == '0001-01-01') {
		payload.DateRemoved = '';
	} else {
		payload.DateRemoved = sp[0];
	}

	sp = payload.DateLifeVows.split('T');
	if (sp[0] == '0001-01-01') {
		payload.DateLifeVows = '';
	} else {
		payload.DateLifeVows = sp[0];
	}
	return payload;
}

export async function updateMember(id, FieldName, Value) {
	const jwt = localStorage.getItem('jwt');
	if (jwt === undefined || jwt === null) {
		throw new Error('Not Logged in');
	}

	const dataArray = new FormData();
	dataArray.append('field', FieldName);
	dataArray.append('value', Value);

	const request = {
		method: 'POST',
		mode: 'cors',
		credentials: 'include',
		redirect: 'manual',
		referrerPolicy: 'origin',
		body: dataArray,
		headers: { Authorization: 'Bearer ' + jwt }
	};

	const response = await fetch(`${server}/api/v1/member/${id}`, request);
	const payload = await response.json();
	if (response.status != 200) {
		console.log('server returned ', response.status);
		throw new Error(payload.error);
	}
}

export async function subsearch(query) {
	const jwt = localStorage.getItem('jwt');
	if (jwt === undefined || jwt === null) {
		throw new Error('Not Logged in');
	}

	const dataArray = new FormData();
	dataArray.append('query', query);

	const request = {
		method: 'POST',
		mode: 'cors',
		credentials: 'include',
		redirect: 'manual',
		referrerPolicy: 'origin',
		body: dataArray,
		headers: { Authorization: 'Bearer ' + jwt }
	};

	const response = await fetch(`${server}/api/v1/subsearch`, request);
	const payload = await response.json();
	if (response.status != 200) {
		console.log('server returned ', response.status);
		throw new Error(payload.error);
	}
	return payload;
}

export async function getSubscriber(id) {
	const jwt = localStorage.getItem('jwt');
	if (jwt === undefined || jwt === null) {
		throw new Error('Not Logged in');
	}

	const request = {
		method: 'GET',
		mode: 'cors',
		credentials: 'include',
		redirect: 'manual',
		referrerPolicy: 'origin',
		headers: { Authorization: 'Bearer ' + jwt }
	};

	const response = await fetch(`${server}/api/v1/subscriber/${id}`, request);
	const payload = await response.json();
	if (response.status != 200) {
		console.log('server returned ', response.status);
		throw new Error(payload.error);
	}

	// trim unnecessary times
	let sp = payload.DatePaid.split('T');
	if (sp[0] == '0001-01-01') {
		payload.DatePaid = '';
	} else {
		payload.DatePaid = sp[0];
	}

	return payload;
}

export async function updateSubscriber(id, FieldName, Value) {
	const jwt = localStorage.getItem('jwt');
	if (jwt === undefined || jwt === null) {
		throw new Error('Not Logged in');
	}

	const dataArray = new FormData();
	dataArray.append('field', FieldName);
	dataArray.append('value', Value);

	const request = {
		method: 'POST',
		mode: 'cors',
		credentials: 'include',
		redirect: 'manual',
		referrerPolicy: 'origin',
		body: dataArray,
		headers: { Authorization: 'Bearer ' + jwt }
	};

	const response = await fetch(`${server}/api/v1/member/${id}`, request);
	const payload = await response.json();
	if (response.status != 200) {
		console.log('server returned ', response.status);
		throw new Error(payload.error);
	}
}

export async function report(reportname) {
	const jwt = localStorage.getItem('jwt');
	if (jwt === undefined || jwt === null) {
		throw new Error('Not Logged in');
	}

	const request = {
		method: 'GET',
		mode: 'cors',
		credentials: 'include',
		redirect: 'manual',
		referrerPolicy: 'origin',
		headers: { Authorization: 'Bearer ' + jwt }
	};

	const response = await fetch(`${server}/api/v1/report/${reportname}`, request);
	const payload = await response.text();
	if (response.status != 200) {
		console.log('server returned ', response.status);
		throw new Error(payload.error);
	}

	// https://stackoverflow.com/questions/14964035/how-to-export-javascript-array-info-to-csv-on-client-side
	const blob = new Blob([payload], { type: 'text/csv;charset=utf-8;' });
	if (navigator.msSaveBlob) {
		// IE 10+
		navigator.msSaveBlob(blob, `${reportname}.csv`);
	} else {
		const link = document.createElement('a');
		if (link.download !== undefined) {
			// feature detection
			// Browsers that support HTML5 download attribute
			const url = URL.createObjectURL(blob);
			link.setAttribute('href', url);
			link.setAttribute('download', `${reportname}.csv`);
			link.style.visibility = 'hidden';
			document.body.appendChild(link);
			link.click();
			document.body.removeChild(link);
		}
	}
}

export async function createMember(firstname, lastname) {
	const jwt = localStorage.getItem('jwt');
	if (jwt === undefined || jwt === null) {
		throw new Error('Not Logged in');
	}

	const dataArray = new FormData();
	dataArray.append('firstname', firstname);
	dataArray.append('lastname', lastname);

	const request = {
		method: 'POST',
		mode: 'cors',
		credentials: 'include',
		redirect: 'manual',
		referrerPolicy: 'origin',
		headers: { Authorization: 'Bearer ' + jwt },
		body: dataArray
	};

	const response = await fetch(`${server}/api/v1/member`, request);
	const payload = await response.json();
	if (response.status != 200) {
		console.log('server returned ', response.status);
		throw new Error(payload.error);
	}
	return payload.id;
}

export async function postRegister(email) {
	const dataArray = new FormData();
	dataArray.append('email', email);

	const request = {
		method: 'POST',
		mode: 'cors',
		credentials: 'include',
		redirect: 'manual',
		referrerPolicy: 'origin',
		body: dataArray
	};

	const response = await fetch(`${server}/api/v1/register`, request);
	const payload = await response.json();
	if (response.status != 200) {
		console.log('server returned ', response.status);
		throw new Error(payload.error);
	}
	return true;
}

export async function getMeFromServer() {
	const jwt = localStorage.getItem('jwt');
	if (jwt === undefined || jwt === null) {
		throw new Error('Not Logged in');
	}
	// const token = JSON.parse(window.atob(jwt.split(".")[1]).toString());

	const request = {
		method: 'GET',
		mode: 'cors',
		credentials: 'include',
		redirect: 'manual',
		referrerPolicy: 'origin',
		headers: { Authorization: 'Bearer ' + jwt }
	};

	const response = await fetch(`${server}/api/v1/me`, request);
	const payload = await response.json();
	if (response.status != 200) {
		console.log('server returned ', response.status);
		throw new Error(payload.error);
	}

	// trim unnecessary times
	let sp = payload.DateReaffirmation.split('T');
	if (sp[0] == '0001-01-01') {
		payload.DateReaffirmation = '';
	} else {
		payload.DateReaffirmation = sp[0];
	}

	sp = payload.DateFirstVows.split('T');
	if (sp[0] == '0001-01-01') {
		payload.DateFirstVows = '';
	} else {
		payload.DateFirstVows = sp[0];
	}

	sp = payload.DateNovitiate.split('T');
	if (sp[0] == '0001-01-01') {
		payload.DateNovitiate = '';
	} else {
		payload.DateNovitiate = sp[0];
	}

	sp = payload.BirthDate.split('T');
	if (sp[0] == '0001-01-01') {
		payload.BirthDate = '';
	} else {
		payload.BirthDate = sp[0];
	}

	sp = payload.DateDeceased.split('T');
	if (sp[0] == '0001-01-01') {
		payload.DateDeceased = '';
	} else {
		payload.DateDeceased = sp[0];
	}

	sp = payload.DateRecordCreated.split('T');
	if (sp[0] == '0001-01-01') {
		payload.DateRecordCreated = '';
	} else {
		payload.DateRecordCreated = sp[0];
	}

	sp = payload.DateRemoved.split('T');
	if (sp[0] == '0001-01-01') {
		payload.DateRemoved = '';
	} else {
		payload.DateRemoved = sp[0];
	}

	sp = payload.DateLifeVows.split('T');
	if (sp[0] == '0001-01-01') {
		payload.DateLifeVows = '';
	} else {
		payload.DateLifeVows = sp[0];
	}
	return payload;
}

export async function updateMe() {
	throw new Error('editing your own data is not yet available');
}

export async function getGiving(id) {
	const jwt = localStorage.getItem('jwt');
	if (jwt === undefined || jwt === null) {
		throw new Error('Not Logged in');
	}

	const request = {
		method: 'GET',
		mode: 'cors',
		credentials: 'include',
		redirect: 'manual',
		referrerPolicy: 'origin',
		headers: { Authorization: 'Bearer ' + jwt }
	};

	const response = await fetch(`${server}/api/v1/giving/${id}`, request);
	const payload = await response.json();
	if (response.status != 200) {
		console.log('server returned ', response.status);
		throw new Error(payload.error);
	}

	payload.forEach((gr) => {
		let sp = gr.Date.split('T');
		gr.Date = sp[0];
	});

	return payload;
}

export async function postGiving(id, date, amount, description, check, transaction) {
	const jwt = localStorage.getItem('jwt');
	if (jwt === undefined || jwt === null) {
		throw new Error('Not Logged in');
	}

	const dataArray = new FormData();
	dataArray.append('id', id);
	dataArray.append('date', date);
	dataArray.append('amount', amount);
	dataArray.append('description', description);
	dataArray.append('check', check);
	dataArray.append('transaction', transaction);

	const request = {
		method: 'POST',
		mode: 'cors',
		credentials: 'include',
		redirect: 'manual',
		referrerPolicy: 'origin',
		body: dataArray,
		headers: { Authorization: 'Bearer ' + jwt }
	};

	const response = await fetch(`${server}/api/v1/giving/${id}`, request);
	const payload = await response.json();
	if (response.status != 200) {
		console.log('server returned ', response.status);
		throw new Error(payload.error);
	}
}
