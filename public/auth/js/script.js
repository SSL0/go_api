let isRegister = false;
let res; // for debugging

// Post request to api
const sendPostRequest = async (url, data, token) => {
    return fetch(url, {
      method: 'POST',
      headers: {'Content-Type': 'application/json',
                'Authorizaion': 'Bearer: ' + token},
      body: JSON.stringify(data)
    })
      .then(response => response.json())
      .catch(error => console.log(error));
}



const switchAuth = () =>{
    if(isRegister){
        document.getElementById('input-container-email').style.display = 'none';
        document.getElementById('switch-auth').textContent = 'Регистрация';
        document.getElementById('btn-submit').textContent = 'Войти в аккаунт';
    } else {
        document.getElementById('input-container-email').style.display = 'block';
        document.getElementById('switch-auth').textContent = 'Вход';
        document.getElementById('btn-submit').textContent = 'Зарегистрироваться';
    }
    isRegister = !isRegister;
}

const submitPress = async () => {
    const name = document.getElementById('input-name').value;
    const email = document.getElementById('input-email').value;
    const pass = document.getElementById('input-pass').value;


    document.getElementById('btn-submit').style.opacity = '0.8';
    if (isRegister && await register(name, email, pass) == false){
        if(!success){
            return;
        }
    }
    await login(name, pass);

    document.getElementById('btn-submit').style.opacity = '';
}

const login = async (name, pass) => {
    const body = {
        name: name,
        password: pass
    }
    const response = await sendPostRequest('/api/auth/login', body);
    if(response.message == "success"){
        res = response;
        document.cookie = "jwt=" + response.token + "; path=/; expires=" + response.expires + "; SameSite=lax";
        console.log(sendPostRequest("/api/user/get-user", body))
        window.location.href = "/home/";
    } else{
        document.getElementById('span-error').style.display = 'block';
        document.getElementById('span-error').textContent = '*' + response.message;
    }

}

const register = async (name, email, pass) => {
    const body = {
        name: name,
        email: email,
        password: pass
    }
    const response = await sendPostRequest('/api/auth/register', body);

    if(response.message == "success"){
        return true;
    } else{
        document.getElementById('span-error').style.display = 'block';
        document.getElementById('span-error').textContent = '*' + response.message;
        return false;
    }
}

const main = async() => {
    // Check, is user unauthorized
    const response = await sendPostRequest("/api/user/get-user", {});
    res = response;
    if(response.message === undefined){
        window.location.href = "/home/";
    } 
 
}

main();