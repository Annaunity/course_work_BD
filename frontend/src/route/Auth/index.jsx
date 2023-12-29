import AppBody from "components/AppBody";
import AppContainer from "components/AppContainer"
import { Button } from 'primereact/button';
import { InputText } from 'primereact/inputtext';
import { Password } from 'primereact/password';
import { Toast } from 'primereact/toast';

import { useState, useRef } from "react";
import { useNavigate } from 'react-router-dom';

import * as api from "api";

import "./index.scoped.css";

function AuthPage(props) {
  const LOGIN_ERROR = "Логин не может быть пустым";
  const PWD_ERROR = "Пароль не может быть пустым";
  const SUCCESS_SIGN_UP = "Вы успешно зарегистрировались";
  const SUCCESS_SIGN_IN = "Вы успешно вошли";

  const toast = useRef();

  const navigate = useNavigate();

  const [usernameUp, setUsernameUp] = useState("");
  const [streetUp, setStreetUp] = useState("");
  const [numberUp, setNumberUp] = useState("");
  const [passwordUp, setPasswordUp] = useState("");
  const [passwordAgainUp, setPasswordAgainUp] = useState("");

  const [usernameIn, setUsernameIn] = useState("");
  const [passwordIn, setPasswordIn] = useState("");

  const [userData, setUserData] = useState(null);

  async function signUp(event) {
    event.preventDefault();

    if (usernameUp === null || usernameUp === "") {
      toast.current.show({ sticky: false, life: 2000, closable: true, severity: "error", summary: "Ошибка", detail: LOGIN_ERROR });
      return;
    }

    if (passwordUp === null || passwordUp === "") {
      toast.current.show({ sticky: false, life: 2000, closable: true, severity: "error", summary: "Ошибка", detail: PWD_ERROR });
      return;
    }

    let data = ({"full_name": usernameUp, "street": streetUp, "house": numberUp, "password": passwordUp});

    let response = await fetch('/api/register', { method: "POST", body: JSON.stringify(data) });
    let message = await response.json();

    console.log(message);
    console.log(response.status);
    if (response.status == 201) {
      if (message.hasOwnProperty("error")) {
        toast.current.show({ sticky: false, life: 2000, closable: true, severity: "error", summary: "Ошибка", detail: message["error"] });
        return;
      }
      // Сохранить данные в storage
      localStorage.setItem("token", message['token']);
      localStorage.setItem("user_id", message['user_id']);
      localStorage.setItem("is_admin", message['is_admin']);
      console.log(localStorage.getItem('token'));
      if (localStorage.getItem('is_admin') === "true") {
        navigate("/Admin");
      } else {
        navigate("/Main");
      }
    } else {
      toast.current.show({
        sticky: false,
        life: 2000,
        closable: true,
        severity: "error",
        summary: "Ошибка",
        detail: message["error"],
        
      });console.log("aboba");
    }
  }


  async function signIn(event) {
    event.preventDefault();

    if (usernameIn === null || usernameIn === "") {
      toast.current.show({ sticky: false, life: 2000, closable: true, severity: "error", summary: "Ошибка", detail: LOGIN_ERROR });
      return;
    }

    if (passwordIn === null || passwordIn === "") {
      toast.current.show({ sticky: false, life: 2000, closable: true, severity: "error", summary: "Ошибка", detail: PWD_ERROR });
      return;
    }

    
    let data = ({"full_name": usernameIn, "password": passwordIn});


    let response = await fetch('/api/auth', { method: "POST", body: JSON.stringify(data) });
    let message = await response.json();

    console.log(message);
    if (response.status == 200) {
      if (message.hasOwnProperty("error")) {
        toast.current.show({ sticky: false, life: 2000, closable: true, severity: "error", summary: "Ошибка", detail: message["error"] });
        return;
      }
      localStorage.setItem("token", message['token']);
      localStorage.setItem("user_id", message['user_id']);
      localStorage.setItem("is_admin", message['is_admin']);
      console.log(localStorage.getItem('token'));
      if (localStorage.getItem('is_admin') === "true") {
        navigate("/Admin");
      } else {
        navigate("/Main");
      }
    }

    else {
      toast.current.show({
        sticky: false,
        life: 2000,
        closable: true,
        severity: "error",
        summary: "Ошибка",
        detail: message["error"],
      });
    }
  }

  return (
    <AppBody>
      <AppContainer>
        <div className="auth">
          <div className="main">
            <input type="checkbox" id="chk" aria-hidden="true"></input>

            <div className="signup">
              <form onSubmit={signUp}>
                <label className="title label-up" htmlFor="chk" aria-hidden="true">Регистрация</label>

                <div className="input-wrapper">
                  <span className="p-float-label">
                    <InputText id="usernameUp" value={usernameUp} onChange={(e) => setUsernameUp(e.target.value)} />
                    <label htmlFor="usernameUp">ФИО</label>
                  </span>

                  <span className="p-float-label">
                    <InputText id="streetUp" value={streetUp} onChange={(e) => setStreetUp(e.target.value)} />
                    <label htmlFor="streetUp">Улица</label>
                  </span>

                  <span className="p-float-label">
                    <InputText id="numberUp" value={numberUp} onChange={(e) => setNumberUp(e.target.value)} />
                    <label htmlFor="numberUp">Дом</label>
                  </span>

                  <span className="p-float-label">
                    <Password inputId="passwordUp" value={passwordUp} onChange={(e) => setPasswordUp(e.target.value)} toggleMask
                      promptLabel="Введите пароль" weakLabel="Слишком легкий" mediumLabel="Средний" strongLabel="Сложный пароль" />
                    <label htmlFor="passwordUp">Пароль</label>
                  </span>

                  <Button className="button-up" label="Создать аккаунт" icon="pi pi-user-plus" iconPos="right" />
                </div>
              </form>
            </div>

            <div className="login">
              <form onSubmit={signIn}>
                <label className="title" htmlFor="chk" aria-hidden="true">Вход</label>
                <div className="input-wrapper">
                  <span className="p-float-label">
                    <InputText id="usernameIn" value={usernameIn} onChange={(e) => setUsernameIn(e.target.value)} />
                    <label htmlFor="usernameIn">ФИО</label>
                  </span>

                  <span className="p-float-label">
                    <Password id="passwordIn" value={passwordIn} onChange={(e) => setPasswordIn(e.target.value)} feedback={false} toggleMask />
                    <label htmlFor="passwordIn">Пароль</label>
                  </span>
                  <Button label="Войти в аккаунт" icon="pi pi-sign-in" iconPos="right" />
                </div>
              </form>
            </div>
          </div>
        </div>
        <div className="footer">

        </div>
      </AppContainer>
      <Toast ref={toast} position="top-right" />
    </AppBody>
  );

}

export default AuthPage;
