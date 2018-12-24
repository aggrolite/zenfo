/* eslint-env browser */
/* global module */
import React from 'react';
import ReactDOM from 'react-dom';

const title = 'zenfo.info';

ReactDOM.render(
  <div>{title}</div>,
  document.getElementById('app')
);

module.hot.accept();
