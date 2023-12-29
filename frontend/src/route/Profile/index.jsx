import React, { useState, useEffect, useRef } from 'react';
import './index.scoped.css';
import Logo from 'route/Image/logo.png';
import { Link, useNavigate } from "react-router-dom";
import { Toast } from 'primereact/toast';

const Profile = () => {
  const [profileData, setProfileData] = useState({
    fullname: '',
    phone: '',
    email: '',
    street: '',
    number: ''
  });



  const [errors, setErrors] = useState({
    nameError: '',
    phoneError: '',
    emailError: '',
  });

  const navigate = useNavigate();
  const toast = useRef();

  useEffect(() => {
    (async () => {
      let response = await fetch('/api/user/' + localStorage.getItem('user_id') + '/profile', {
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
      setProfileData(message);
    })()
  }, [])


  const handleSave = async (event) => {
    event.preventDefault();
    // const isFormValid = checkFormValidity();

    let data = { "full_name": profileData.fullname, "street": profileData.street, "house": profileData.number, "email": profileData.email, "phone": profileData.phone };

    let response = await fetch('/api/user/' + localStorage.getItem('user_id') + '/profile', {
      method: 'PATCH',
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
      detail: 'Профиль обновлен'
    });
    setProfileData(message);
  };


  const handleChange = (e) => {
    const { name, value } = e.target;
    setProfileData({ ...profileData, [name]: value });
  };


  const checkFormValidity = () => {
    const { name, phone, email } = profileData;
    const errors = {};

    if (!name.match(/^[a-zA-Zа-яА-Я]+$/)) {
      errors.nameError = 'Имя может содержать только буквы';
    }

    if (!phone.match(/^\d{11}$/)) {
      errors.phoneError = 'Номер телефона должен состоять из 11 цифр';
    }

    if (!email.match(/^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,3})+$/)) {
      errors.emailError = 'Некорректный адрес электронной почты';
    }

    setErrors(errors);

    return Object.keys(errors).length === 0;
  };

  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('is_admin');
    localStorage.removeItem('user_id');
    navigate("/Auth")
    alert('Вы вышли из аккаунта');
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
          <li><Link to="/Main">Главная</Link></li>
          <li><Link to="/Profile">Профиль</Link></li>
          <li><Link to="/Contacts">Контакты</Link></li>
        </ul>
      </nav>

      <div className="content">
        <div className="profile-card">
          <h1>Личный профиль</h1>
          <div className="row">
            <label htmlFor="name">ФИО:</label>
            <input
              type="text"
              id="name"
              name="name"
              value={profileData.fullname}
              onChange={handleChange}
            />
            {errors.nameError && <span className="error">{errors.nameError}</span>}
          </div>
          <div className="row">
            <label htmlFor="phone">Телефон:</label>
            <input
              type="text"
              id="phone"
              name="phone"
              value={profileData.phone}
              onChange={handleChange}
            />
            {errors.phoneError && <span className="error">{errors.phoneError}</span>}
          </div>
          <div className="row">
            <label htmlFor="email">E-mail:</label>
            <input
              type="text"
              id="email"
              name="email"
              value={profileData.email}
              onChange={handleChange}
            />
            {errors.emailError && <span className="error">{errors.emailError}</span>}
          </div>
          <div className="row">
            <label htmlFor="street">Улица:</label>
            <input
              type="text"
              id="street"
              name="street"
              value={profileData.street}
              onChange={handleChange}
              readOnly
            />
          </div>
          <div className="row">
            <label htmlFor="house">Дом:</label>
            <input
              type="text"
              id="house"
              name="house"
              value={profileData.number}
              onChange={handleChange}
              readOnly
            />
          </div>
          <div className="button-row">
            <button className="save-button" onClick={handleSave}>
              Сохранить
            </button>
            <button className="logout-button" onClick={handleLogout} >
              Выход
            </button>
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

export default Profile;