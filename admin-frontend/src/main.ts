import { createApp } from 'vue';
import App from './App.vue';
import router from './router';
import { i18n } from './i18n';
import './style.css';

createApp(App).use(router).use(i18n).mount('#app');

console.log(
  `%c
  ...     _M_
 /( )\\    ( )
/ / \\ \\  / : \\
~~\\%/~~  \\|:|/
 /   \\    |||
/,,,,,\\   |||     Wedding Card

Not sure what you're doing here, but thanks!
`,
  'font-family: monospace;',
);
