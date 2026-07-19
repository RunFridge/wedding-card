import { createApp } from 'vue';
import App from './App.vue';
import router from './router';
import { i18n } from './i18n';
import './style.css';
import { loadConfig } from './config/wedding';

loadConfig().then(() => {
  createApp(App).use(router).use(i18n).mount('#app');

  console.log(
    `%c
   .mPMMNNHHmo.       _, n!it6XHm.
 ,8PMWWHHKKDDXY8.   _o86SSXXDDKKH8b.
,8FMWNKKQDXX665Y8. d8YJYY5566XXDQKY8.
d8MWNKKDDSS55YYtY88PjjtjJtYY55SSDDK8b
Y8NNKKDXS65YJtjjiPPi=i=iijjtJY56SXDY8
i8WHKDDS65Yttcc==++>+>++==ccttY56SDd8
\`8bKQDSS5Yttci=+>!;;:;;!>+=icttY5SSd8
 Y8KQXX65JJjc==>!::~~~::!>==cjJJ56X8P
 \`8LDX66YJjji=>>;:'. .':;>>=ijjJY668'
  i8QXX65JJjc==>!::~~~::!>==cjJJ568P
   Y8DSS5Yttci=+>!;;:;;!>+=icttY58P
    8bDS65Yttcc==++>+>++==ccttY58P
     Y8XS65YJtjjii=i=i=iijjtJY58'
      \`8bSS55YYtJjtjjjtjJtYY5d8'
        \`8oX6655YYJYJYJYY556dP
          \`8oXXSS6666666SSo8'
            \`YbDXDXXXXXo8P'
              \`YbKKQKdP'
                \`Yb8P'
                  Y8
                   ' mh

 You found me! check out https://github.com/RunFridge/wedding-card
`,
    'color: #e91e63; font-family: monospace;',
  );
});
