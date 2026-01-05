import 'vuetify/styles'
import { createVuetify } from 'vuetify'
import { aliases, mdi } from 'vuetify/iconsets/mdi-svg'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { ru } from 'vuetify/locale'

export default createVuetify({
  components,
  directives,
  icons: {
    defaultSet: 'mdi',
    aliases,
    sets: { mdi }
  },
  locale: {
    locale: 'ru',
    messages: { ru }
  },
  theme: {
    defaultTheme: 'light',
    themes: {
      light: {
        dark: false,
        colors: {
          background: '#f8fafc',
          surface: '#ffffff',
          'surface-variant': '#f1f5f9',
          primary: '#1e293b',
          'primary-darken-1': '#0f172a',
          secondary: '#64748b',
          'secondary-darken-1': '#475569',
          error: '#dc2626',
          warning: '#d97706',
          info: '#0284c7',
          success: '#059669',
          'on-background': '#1e293b',
          'on-surface': '#334155',
          'on-primary': '#ffffff',
          'on-secondary': '#ffffff'
        }
      }
    }
  },
  defaults: {
    VBtn: {
      variant: 'flat',
      rounded: 'lg'
    },
    VCard: {
      rounded: 'lg',
      elevation: 1
    },
    VTextField: {
      variant: 'outlined',
      density: 'comfortable',
      rounded: 'lg'
    },
    VSelect: {
      variant: 'outlined',
      density: 'comfortable',
      rounded: 'lg'
    },
    VCheckbox: {
      density: 'comfortable'
    },
    VDataTable: {
      density: 'comfortable'
    },
    VDialog: {
      rounded: 'lg'
    },
    VChip: {
      rounded: 'lg'
    },
    VAlert: {
      rounded: 'lg'
    }
  }
})
