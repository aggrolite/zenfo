/* eslint-env browser */
/* global module */
import React from 'react';
import ReactDOM from 'react-dom';
import App from './component/App';
import './style.css';

ReactDOM.render(
  <App />,
  document.getElementById('app')
);

module.hot.accept();
