export const server = "http://saint-luke.net:8443";

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
  dataArray.append("username", username);
  dataArray.append("password", password);

  const request = {
    method: "POST",
    mode: "cors",
    credentials: "include",
    redirect: "manual",
    referrerPolicy: "origin",
    body: dataArray,
  };

  const response = await fetch(`${server}/api/v1/getJWT`, request);
  const payload = await response.text();
  if (response.status != 200) {
    console.log("server returned ", response.status);
    throw new Error(payload.error);
  }
  localStorage.setItem("jwt", payload);
}

export async function search(query) {
  const jwt = localStorage.getItem('jwt');
  if (jwt === undefined || jwt === null) {
    throw new Error("Not Logged in");
  }

  const dataArray = new FormData();
  dataArray.append("query", query);

  const request = {
    method: "POST",
    mode: "cors",
    credentials: "include",
    redirect: "manual",
    referrerPolicy: "origin",
    body: dataArray,
    headers: {Authorization: "Bearer " + jwt}
  };

  const response = await fetch(`${server}/api/v1/search`, request);
  const payload = await response.json();
  if (response.status != 200) {
    console.log("server returned ", response.status);
    throw new Error(payload.error);
  }
  return payload;
}

export async function getMember(id) {
  const jwt = localStorage.getItem('jwt');
  if (jwt === undefined || jwt === null) {
    throw new Error("Not Logged in");
  }

  const request = {
    method: "GET",
    mode: "cors",
    credentials: "include",
    redirect: "manual",
    referrerPolicy: "origin",
    headers: {Authorization: "Bearer " + jwt}
  };

  const response = await fetch(`${server}/api/v1/member/${id}`, request);
  const payload = await response.json();
  if (response.status != 200) {
    console.log("server returned ", response.status);
    throw new Error(payload.error);
  }

  // trim unnecessary times
  let sp = payload.DateReaffirmation.split('T');
  if ( sp[0] == "1800-01-01" ) {
    payload.DateReaffirmation = '';
  } else {
    payload.DateReaffirmation = sp[0];
  }

  sp = payload.DateFirstProfession.split('T');
  if ( sp[0] == "1800-01-01" ) {
    payload.DateFirstProfession= '';
  } else {
    payload.DateFirstProfession = sp[0];
  }

  sp = payload.DateFirstVows.split('T');
  if ( sp[0] == "1800-01-01" ) {
    payload.DateFirstVows= '';
  } else {
    payload.DateFirstVows = sp[0];
  }
  
  sp = payload.DateNovitiate.split('T');
  if ( sp[0] == "1800-01-01" ) {
    payload.DateNovitiate= '';
  } else {
    payload.DateNovitiate = sp[0];
  }

  sp = payload.BirthDate.split('T');
  if ( sp[0] == "1800-01-01" ) {
    payload.BirthDate= '';
  } else {
    payload.BirthDate = sp[0];
  }

  sp = payload.DateDeceased.split('T');
  if ( sp[0] == "1800-01-01" ) {
    payload.DateDeceased= '';
  } else {
    payload.DateDeceased = sp[0];
  }

  sp = payload.DateRecordCreated.split('T');
  if ( sp[0] == "1800-01-01" ) {
    payload.DateRecordCreated= '';
  } else {
    payload.DateRecordCreated = sp[0];
  }

  sp = payload.DateRemoved.split('T');
  if ( sp[0] == "1800-01-01" ) {
    payload.DateRemoved = '';
  } else {
    payload.DateRemoved = sp[0];
  }

  return payload;
}

export async function updateMember(id, FieldName, Value) {
  const jwt = localStorage.getItem('jwt');
  if (jwt === undefined || jwt === null) {
    throw new Error("Not Logged in");
  }

  const dataArray = new FormData();
  dataArray.append("field", FieldName);
  dataArray.append("value", Value);

  const request = {
    method: "POST",
    mode: "cors",
    credentials: "include",
    redirect: "manual",
    referrerPolicy: "origin",
    body: dataArray,
    headers: {Authorization: "Bearer " + jwt}
  };

  const response = await fetch(`${server}/api/v1/member/${id}`, request);
  const payload = await response.json();
  if (response.status != 200) {
    console.log("server returned ", response.status);
    throw new Error(payload.error);
  }
}

export async function subsearch(query) {
  const jwt = localStorage.getItem('jwt');
  if (jwt === undefined || jwt === null) {
    throw new Error("Not Logged in");
  }

  const dataArray = new FormData();
  dataArray.append("query", query);

  const request = {
    method: "POST",
    mode: "cors",
    credentials: "include",
    redirect: "manual",
    referrerPolicy: "origin",
    body: dataArray,
    headers: {Authorization: "Bearer " + jwt}
  };

  const response = await fetch(`${server}/api/v1/subsearch`, request);
  const payload = await response.json();
  if (response.status != 200) {
    console.log("server returned ", response.status);
    throw new Error(payload.error);
  }
  return payload;
}

export async function getSubscriber(id) {
  const jwt = localStorage.getItem('jwt');
  if (jwt === undefined || jwt === null) {
    throw new Error("Not Logged in");
  }

  const request = {
    method: "GET",
    mode: "cors",
    credentials: "include",
    redirect: "manual",
    referrerPolicy: "origin",
    headers: {Authorization: "Bearer " + jwt}
  };

  const response = await fetch(`${server}/api/v1/subscriber/${id}`, request);
  const payload = await response.json();
  if (response.status != 200) {
    console.log("server returned ", response.status);
    throw new Error(payload.error);
  }

  // trim unnecessary times
  let sp = payload.DatePaid.split('T');
  if ( sp[0] == "1800-01-01" ) {
    payload.DatePaid= '';
  } else {
    payload.DatePaid = sp[0];
  }

  return payload;
}

export async function updateSubscriber(id, FieldName, Value) {
  const jwt = localStorage.getItem('jwt');
  if (jwt === undefined || jwt === null) {
    throw new Error("Not Logged in");
  }

  const dataArray = new FormData();
  dataArray.append("field", FieldName);
  dataArray.append("value", Value);

  const request = {
    method: "POST",
    mode: "cors",
    credentials: "include",
    redirect: "manual",
    referrerPolicy: "origin",
    body: dataArray,
    headers: {Authorization: "Bearer " + jwt}
  };

  const response = await fetch(`${server}/api/v1/member/${id}`, request);
  const payload = await response.json();
  if (response.status != 200) {
    console.log("server returned ", response.status);
    throw new Error(payload.error);
  }
}
