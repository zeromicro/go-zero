import i18n from 'i18next';
import {initReactI18next} from 'react-i18next';
import {resources} from './locales/resources';

i18n
    .use(initReactI18next)
    .init({
        resources,
        fallbackLng: "en",
        lng: 'zh',
        debug: true,
        interpolation: {
            escapeValue: false,
        }
    });

export default i18n;
