import React from 'react';
import { Link, useNavigate } from "react-router-dom";
import "./index.scoped.css";
import Logo from 'route/Image/logo.png';
import SimpleLineIcon from 'react-simple-line-icons';

function Main() {
  const navigate = useNavigate();

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
      <div className="content">
        <h1>О нас</h1>
        <p>
          Работаем с 2023 года. Оказываем услуги капитального ремонта и приема жалоб и предложений.
        </p>
        <div className="card-container">
          <div className="card">
            <i className="icon-user success font-large-2 float-right"></i>
            <h2><Link to="/Profile">Личный профиль</Link></h2>
            <p>Посмотрите свои личные данные</p>
          </div>
          <div className="card">
            <i className="icon-pencil primary font-large-2 float-left"></i>
            <h2><Link to="/Complain">Куда жаловаться</Link></h2>
            <p>Отправьте свой запрос в компанию</p>
          </div>
          <div className="card">
            <i className="icon-graph success font-large-2 float-left"></i>
            <h2><Link to="/Stat">Статистика дома</Link></h2>
            <p>Статистика и отзывы вашего дома</p>
          </div>
          <div className="card">
            <i className="icon-speech warning font-large-2 float-left"></i>
            <h2><Link to="/Status">Статус жалобы</Link></h2>
            <p>Проверка текущего статуса вашей жалобы</p>
          </div>
        </div>
      </div>
      <footer className="footer">
        <p>Тел. горячей линии: +7 (800) 555-35-35</p>
        <p>ООО "Наше предприятие", 2023 г.</p>
        <p>E-mail: info@gmail.com</p>
      </footer>
    </div>
  );
}

export default Main;