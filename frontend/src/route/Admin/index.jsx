import React, { useState, useEffect, useRef } from 'react';
import './index.scoped.css';
import Logo from 'route/Image/logo.png';
import { Link, useNavigate } from "react-router-dom";
import { Toast } from 'primereact/toast';


const Admin = () => {
  const [isRatingEntered, setIsRatingEntered] = useState(false);
  const [rating, setRating] = useState('');
  const [isButtonDisabled, setIsButtonDisabled] = useState(true);
  const [complaintData, setComplaintData] = useState({});

  const navigate = useNavigate();
  const toast = useRef();

  const handleRatingChange = (event) => {
    setRating(event.target.value);
    setIsRatingEntered(true);
  };

  async function handleAccept(event) {
    event.preventDefault();
    if (isRatingEntered) {
      let data = ({ "complaint_id": complaintData.complaint_id.toString(), "verdict": "принято", "rating": parseFloat(rating)});

      let response = await fetch('/api/admin/' + localStorage.getItem('user_id') + '/report', {
        method: "POST",
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + localStorage.getItem('token')
        },
        body: JSON.stringify(data)
      });

      let message = await response.json();

      console.log(message);
      if (response.status == 200) {
        if (message.hasOwnProperty("error")) {
          toast.current.show({ sticky: false, life: 2000, closable: true, severity: "error", summary: "Ошибка", detail: message["error"] });
          return;
        }
        setComplaintData(message);
      } else {
        toast.current.show({
          sticky: false,
          life: 2000,
          closable: true,
          severity: "error",
          summary: "Ошибка",
          detail: message["error"],
        });
      }

    } else {
    }
  };

  async function handleReject(event) {
    event.preventDefault();
    if (isRatingEntered) {
      let data = ({ "complaint_id": complaintData.complaint_id.toString(), "verdict": "отклонено", "rating": parseFloat(rating) });

      let response = await fetch('/api/admin/' + localStorage.getItem('user_id') + '/report', {
        method: "POST",
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + localStorage.getItem('token')
        },
        body: JSON.stringify(data)
      });

      let message = await response.json();

      console.log(message);
      if (response.status == 200) {
        if (message.hasOwnProperty("error")) {
          toast.current.show({ sticky: false, life: 2000, closable: true, severity: "error", summary: "Ошибка", detail: message["error"] });
          return;
        }
        setComplaintData(message);
      } else {
        toast.current.show({
          sticky: false,
          life: 2000,
          closable: true,
          severity: "error",
          summary: "Ошибка",
          detail: message["error"],
        });
      }

    } else {
    }

  };

  useEffect(() => {
    (async () => {
      if (localStorage.getItem('is_admin') == 'false') {
        navigate("/Auth")
      }
      try {
      let response = await fetch('/api/admin/' + localStorage.getItem('user_id') + '/report', {
        method: "GET",
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + localStorage.getItem('token')
        }
      });
    
      let message = await response.json();
    
      console.log(response.status);
      if (response.status == 401) {
        console.log('aboba');
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
      setComplaintData(message);
    } catch (e) {navigate("/Auth")}
    })()
  }, [])

  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('is_admin');
    localStorage.removeItem('user_id');
    navigate("/Auth")
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
          <span>&#x260E;</span>
          <p>+7 (800) 555-35-35</p>
        </div>
      </header>
      <nav className="nav">
        <ul>
          <li onClick={handleLogout}>Выход</li>
        </ul>
      </nav>
  
      <div className="admin-page">
        <div className="card">
          <div className="card-content">
            <h2>Автор жалобы: {complaintData.full_name}</h2>
            <h3>Тема жалобы: {complaintData.title}</h3>
            <p>Описание жалобы: {complaintData.text}</p>
            <p>Улица: {complaintData.street}</p>
            <p>Дом: {complaintData.number_of_house}</p>
            <p>Год постройки: {complaintData.year_construct}</p>
            <p>Рейтинг: {complaintData.evaluation}</p>
            <p>Дата отправки: {complaintData.date}</p>
            <div className="rating-container">
              <label htmlFor="rating">Оцените жалобу:</label>
              <input
                type="number"
                id="rating"
                min="1"
                max="5"
                value={rating}
                onChange={handleRatingChange}
              />
            </div>
            <div className="button-container">
              <button
                onClick={handleAccept}
                disabled={!isRatingEntered}
                className={`yes ${!isRatingEntered ? 'disabled' : ''}`}
              >
                Принято
              </button>
              <button
                onClick={handleReject}
                disabled={!isRatingEntered}
                className={`no ${!isRatingEntered ? 'disabled' : ''}`}
              >
                Отклонено
              </button>
            </div>
          </div></div>
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

export default Admin;