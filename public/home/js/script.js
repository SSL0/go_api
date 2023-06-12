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

const logout = () => {
  document.cookie = "jwt= ; path=/; expires=" + Date.now(); + "; SameSite=Lax; HttpOnly";
  window.location.href = "/auth/";
}


const main = async() => {
    // Check, is user authorized
    const response = await sendPostRequest("/api/auth/get-user", {});
    res = response;
    if(response.message == 'unauthenticated'){
        window.location.href = "/auth/";
    }
    document.getElementById('div-greeting').textContent = "Привет, " + response.name + "!";
        
}

main();