import React, { useState, useEffect, useRef } from 'react';
import './index.scoped.css';
import Logo from 'route/Image/logo.png';
import { Toast } from 'primereact/toast';
import { Link, useNavigate } from "react-router-dom";

const ComplaintPage = () => {

  const [StatusData, setStatusData] = useState([]);
  const navigate = useNavigate();
  const toast = useRef();

  useEffect(() => {
    (async () => {
      let response = await fetch('/api/user/' + localStorage.getItem('user_id') + '/report', {
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
      setStatusData(message);
    })()
  }, [])

  const getStatusColor = (complaint_status) => {
    switch (complaint_status) {
      case 'принято':
        return 'green';
      case 'в обработке':
        return 'yellow';
      case 'отклонено':
        return 'red';
      default:
        return '';
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

      <div className="status-page">
        <div className="card">
          <table>
            <thead>
              <tr>
                <th>Тема жалобы</th>
                <th>Дата отправки</th>
                <th>Статус</th>
              </tr>
            </thead>
            <tbody>
              {StatusData.map((statusData, index) => (
                <tr key={index}>
                  <td>{statusData.title}</td>
                  <td>{statusData.complaint_data}</td>
                  <td className={`status-indicator ${getStatusColor(statusData.complaint_status)}`}>
                    {statusData.complaint_status}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
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

export default ComplaintPage;