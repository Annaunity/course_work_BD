import React, { useState, useEffect, useRef } from 'react';
import './index.scoped.css';
import Logo from 'route/Image/logo.png';
import { Link, useNavigate } from "react-router-dom";
import { Toast } from 'primereact/toast';

function Contacts() {
  const [companyName, setCompanyName] = useState({company_name:''});
  const [phoneNumber, setPhoneNumber] = useState('');


  const navigate = useNavigate(); 
  const toast = useRef();

  useEffect(() => { 
    (async () => {
      let response = await fetch('/api/user/' + localStorage.getItem('user_id') + '/support', { method: "GET", 
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': 'Bearer ' + localStorage.getItem('token') 
      }});
      let message = await response.json();
      console.log(message);
      if (response.status == 401) {
        localStorage.removeItem('token');
        localStorage.removeItem('is_admin');
        localStorage.removeItem('user_id');
        toast.current.show({
          sticky: false,
          life: 2000,
          closable: true,
          severity: "error",
          summary: "Ошибка",
          detail: message["error"],
        });
        await new Promise(resolve => setTimeout(resolve, 500));
        navigate("/Auth")
      } 
      if (message.hasOwnProperty("error")) {
        toast.current.show({ sticky: false, life: 2000, closable: true, severity: "error", summary: "Ошибка", detail: message["error"] });
        return;
      } 
      setCompanyName(message);
      setPhoneNumber(message.phone);
    })()
  }, []) 

  return (
    <div className="App">
      <header className="header">
        <div className="logo">
          <img src={Logo} alt="logo" />
          <h4>Капитальные работы по г. Санкт-Петербург</h4>
        </div>
        <div className="hotline">
          <p>Тел. горячей линии:</p>
          <span>&#x260E; </span>
          <p>+7 (800) 555-35-35</p>
        </div>
      </header>
      <nav className="nav">
      <ul>
          <li><Link to="/Main">Главная</Link></li>
          <li><Link to="/Profile">Профиль</Link></li>
          <li><Link to="/Contacts">Контакты</Link></li>
        </ul>
      </nav>

      <div className="company-card">
        <h3>Ваша управляющая компания</h3>
        <p>{companyName.company_name}</p>
        <div className='phone'>
        <h3><lable>Телефон</lable></h3>
        <p>{phoneNumber}</p>
        </div>
      </div>

      <footer className="footer">
        <p>Тел. горячей линии: +7 (800) 555-35-35</p>
        <p>ООО "Наше предприятие", 2023 г.</p>
        <p>E-mail: info@gmail.com</p>
      </footer>
      <Toast ref={toast} position="top-right" />
    </div>
  );
}

export default Contacts;