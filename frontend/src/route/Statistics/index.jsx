import React, { useState, useEffect, useRef } from 'react';
import './index.scoped.css';
import Logo from 'route/Image/logo.png';
import { Link, useNavigate } from "react-router-dom";
import { Toast } from 'primereact/toast';
import House from 'route/Image/house.jpg';

const Statistics = () => {
  const [reviews, setReviews] = useState([]);
  const [newReview, setNewReview] = useState('');

  const [houseData, setHouseData] = useState({
    street_house: '',
    number_house: '',
    year_construct: '',
    number_of_apartments: '',
    evaluation: '',
    discounts: ''
  });

  const [houseFeed, setHouseFeed] = useState({
    user_fullname: '',
    feedback_text: '',
    feedback_date: ''
  });

  const navigate = useNavigate();
  const toast = useRef();


  useEffect(() => {
    (async () => {
      let response = await fetch('/api/user/' + localStorage.getItem('user_id') + '/house', {
        method: "GET",
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + localStorage.getItem('token')
        }
      });
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
      setHouseData(message);

      let response2 = await fetch('/api/user/' + localStorage.getItem('user_id') + '/feedback', {
        method: "GET",
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + localStorage.getItem('token')
        }
      });
      let message2 = await response2.json();
      console.log(message2);
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
          detail: message2["error"],
        });
        await new Promise(resolve => setTimeout(resolve, 500));
        navigate("/Auth")
      }
      if (message2.hasOwnProperty("error")) {
        toast.current.show({ sticky: false, life: 2000, closable: true, severity: "error", summary: "Ошибка", detail: message2["error"] });
        return;
      }
      setReviews(message2 );
    })()
  }, [])

  async function handleSubmit(event) {
    event.preventDefault();

    let data = {
      feedback_text: newReview
    };
    
    let response = await fetch('/api/user/' + localStorage.getItem('user_id') + '/feedback', {
      method: 'POST',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': 'Bearer ' + localStorage.getItem('token')
      },
      body: JSON.stringify(data)
    });
  
    let message = await response.json();
    console.log(message);
    if (response.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('is_admin');
      localStorage.removeItem('user_id');
      toast.current.show({
        sticky: false,
        life: 2000,
        closable: true,
        severity: 'error',
        summary: 'Ошибка',
        detail: message['error']
      });
      await new Promise(resolve => setTimeout(resolve, 500));
      navigate('/Auth');
    }
    if (message.hasOwnProperty('error')) {
      toast.current.show({
        sticky: false,
        life: 2000,
        closable: true,
        severity: 'error',
        summary: 'Ошибка',
        detail: message['error']
      });
      return;
    }
    toast.current.show({
      sticky: false,
      life: 2000,
      closable: true,
      severity: 'success',
      summary: 'Успешно',
      detail: 'Обратная связь отправлена'
    });
    setReviews(message);
    setNewReview('');
  };

  const handleInputChange = (event) => {
    setNewReview(event.target.value);
  };

  if (!houseData) {
    return <div>Loading...</div>;
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
          <li><Link to="/Main">Главная</Link></li>
          <li><Link to="/Profile">Профиль</Link></li>
          <li><Link to="/Contacts">Контакты</Link></li>
        </ul>
      </nav>

      <div className="container">
        <div className="glass-card">
          <div className="card">
            <div className="card__image">
              <img src={House} alt="House" />
            </div>
            <div className="card__info">
              <h2>
                Адрес: {houseData.street_house} {houseData.number_house}
              </h2>
              <p>Год постройки: {houseData.year_construct}</p>
              <p>Количество этажей: {houseData.num_of_floors}</p>
              <p>Количество квартир: {houseData.number_of_apartments}</p>
              <p>Рейтинг дома: {houseData.evaluation}</p>
              <p>Скидка: {houseData.discount}</p>
            </div>
          </div>

          <h3>Оставить отзыв</h3>
          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <textarea
                id="review"
                name="review"
                value={newReview}
                onChange={handleInputChange}
                required
                placeholder="Напишите отзыв"
              />
            </div>
            <button type="submit" className="button" onClick={handleSubmit}>
              Отправить
            </button>
          </form>

          <div className="reviews">
            {reviews.map((reviews, index) => (
              <tr key={index}>
                <td>{reviews.feedback_text}</td>
              </tr>
            ))}
          </div>
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
};

export default Statistics;