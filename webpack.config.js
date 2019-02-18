const webpack = require('webpack');
const path = require('path');
const Copy = require('copy-webpack-plugin');

module.exports = {
  entry: './src/index.js',
  module: {
    rules: [
      {
        test: /\.(js|jsx)$/,
        exclude: /node_modules/,
        use: ['babel-loader', 'eslint-loader']
      },
      {
        test: /\.css$/,
        use: ['style-loader', 'css-loader'],
      },
    ]
  },
  resolve: {
    extensions: ['*', '.js', '.jsx']
  },
  output: {
    path: __dirname + '/dist',
    publicPath: '/',
    filename: 'bundle.js'
  },
  plugins: [
    new webpack.HotModuleReplacementPlugin(),
    new Copy([
      // relative path is from src
      { from: './static/favicon.ico' }, // <- your path to favicon
    ]),
  ],
  devServer: {
    contentBase: path.join(__dirname, 'dist'),
    compress: true, port: 8080,
    historyApiFallback: true,
    publicPath: "/"
  }
};
