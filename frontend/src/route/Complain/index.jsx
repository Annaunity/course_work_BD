import React, { useState, useRef } from 'react';
import './index.scoped.css';
import Logo from 'route/Image/logo.png';
import { Link, useNavigate } from 'react-router-dom';
import { Toast } from 'primereact/toast';

const ComplaintPage = () => {
  const [complaintTitle, setComplaintTitle] = useState('');
  const [complaintText, setComplaintText] = useState('');
  const [error, setError] = useState(false);
  const toast = useRef(null);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!complaintTitle || !complaintText) {
      setError(true);
    } else {
      setError(false);
      await feedback();
    }
  };

  const feedback = async () => {
    let data = {
      title: complaintTitle,
      text: complaintText,
    };

    let response = await fetch('/api/user/' + localStorage.getItem('user_id') + '/report', {
      method: 'POST',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': 'Bearer ' + localStorage.getItem('token')
      },
      body: JSON.stringify(data),
    });

    let message = await response.json();

    console.log(message);
    if (response.status === 200) {
      if (message.hasOwnProperty('error')) {
        toast.current.show({
          sticky: false,
          life: 2000,
          closable: true,
          severity: 'error',
          summary: 'Ошибка',
          detail: message['error'],
        });
        return;
      }
      navigate('/Main'); 
    }
  };

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
          <li>
            <Link to="/Main">Главная</Link>
          </li>
          <li>
            <Link to="/Profile">Профиль</Link>
          </li>
          <li>
            <Link to="/Contacts">Контакты</Link>
          </li>
        </ul>
      </nav>

      <div className="complaint-page">
        <form className="card" onSubmit={handleSubmit}>
          <h3>Тема:</h3>
          <div className="row">
            <input
              type="text"
              id="complaint-title"
              value={complaintTitle}
              onChange={(e) => setComplaintTitle(e.target.value)}
              placeholder="Напишите тему обращения"
            />
          </div>

          <h3>Описание:</h3>
          <div className="row">
            <textarea
              id="complaint-text"
              value={complaintText}
              onChange={(e) => setComplaintText(e.target.value)}
              placeholder="Напишите описание"
            />
          </div>

          {error && <p>Все поля должны быть заполнены</p>}
          <div className="row">
            <button type="submit" className="button">
              Отправить
            </button>
          </div>
        </form>
      </div>

      <footer className="footer">
        <p>Тел. горячей линии: +7 (800) 555-35-35</p>
        <p>ООО "Наше предприятие", 2023 г.</p>
        <p>E-mail: info@gmail.com</p>
      </footer>
    </div>
  );
};

export default ComplaintPage;