let res; // for debugging

const sendPostRequest = (url, data) => {
    return fetch(url, {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify(data)
    })
      .then(response => response.json())
      .catch(error => console.log(error));
}
const sendGetRequest = (url) => {
  return fetch(url)
    .then(response => response.json())
    .catch(error => console.log(error));
}

const logoutClicked = () => {
  document.cookie = "jwt= ; path=/; expires=" + Date.now(); + "; SameSite=Lax; HttpOnly";
  window.location.href = "/auth/";
}

const clickButtonClicked = async () => {
  let response = await sendGetRequest("/api/user/click");
  response = await sendGetRequest("/api/user/balance");
  document.getElementById('balance-value').textContent = response.balance;
}

const main = async() => {
    // Check, is user authorized
    let response = await sendPostRequest("/api/user/get-user", {});
    res = response;
    if(response.message == 'unauthenticated'){
        window.location.href = "/auth/";
    }
    document.getElementById('div-greeting').textContent = "Привет, " + response.name + "!";
    response = await sendGetRequest("/api/user/balance");
    document.getElementById('balance-value').textContent = response.balance;
        
}

main();