'use strict'
const { VueLoaderPlugin } = require('vue-loader');
const path = require('path');

module.exports = {
    mode: 'development',
    entry: {
        base: './resources/views/js/base.js',
        vue: './resources/views/js/vue.js',
        material: './resources/views/js/material.js',
    },
    output: {
        path: path.join(__dirname, './resources/assets/js'),
        filename: "[name].js"
    },
    module: {
        rules: [
            {
                test: /\.vue$/,
                use: 'vue-loader'
            },
            {
                test: /\.js$/,
                loader: 'babel-loader',
                include: [path.join(__dirname, 'resources/views/js')],
            },
            {
                test: /\.scss$/,
                use: [
                    'vue-style-loader',
                    'css-loader',
                    {
                        loader: 'sass-loader',
                    },
                ],
            },
            {
                test: /\.css$/,
                use: [
                    'handlebars-loader', // handlebars loader expects raw resource string
                    'extract-loader',
                    'css-loader',
                ]
            },
            {
                test: /\.sass$/,
                use: [
                    'sass-loader',
                ]
            }
        ]
    },
    plugins: [
        new VueLoaderPlugin()
    ]
};