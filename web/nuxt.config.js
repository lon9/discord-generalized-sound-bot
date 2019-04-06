module.exports = {
  /*
  ** Headers of the page
  */
  head: {
    title: process.env.WEB_TITLE || 'Discord voice bot',
    meta: [
      { charset: 'utf-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      { hid: 'description', name: 'description', content:  process.env.WEB_DESCRIPTION || 'Discord voice bot'},
      { hid: 'og:type', property: 'og:type', content: 'website'},
      { hid: 'og:title', property: 'og:title', content: process.env.WEB_TITLE},
      { hid: 'og:url', property: 'og:url', content: process.env.WEB_OG_URL},
      { hid: 'og:image', property: 'og:image', content: process.env.WEB_OG_IMAGE},
      { hid: 'og:description', property: 'og:description', content: process.env.WEB_OG_DESCRIPTION},
      { hid: 'twitter:card', name: 'twitter:card', content: process.env.WEB_TWITTER_CARD},
      { hid: 'twitter:title', name: 'twitter:title', content: process.env.WEB_TITLE},
      { hid: 'twitter:description', name: 'twitter:description', content: process.env.WEB_OG_DESCRIPTION},
      { hid: 'twitter:image', name: 'twitter:image', content: process.env.WEB_OG_IMAGE}
    ],
    link: [
      { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' }
    ]
  },
  /*
  ** Customize the progress bar color
  */
  loading: { color: '#3B8070' },
  /*
  ** Build configuration
  */
  build: {
    /*
    ** Run ESLint on save
    */
    extend (config) {
      if (process.server && process.browser) {
        config.module.rules.push({
          enforce: 'pre',
          test: /\.(js|vue)$/,
          loader: 'eslint-loader',
          exclude: /(node_modules)/
        })
      }
    }
  },
  plugins: [
    '~plugins/element-ui',
    '~plugins/clipboard'
  ],
  css: [
    'element-ui/lib/theme-chalk/index.css'
  ],
  env: {
    baseUrl: process.env.BASE_URL || 'http://localhost:3000',
    title: process.env.WEB_TITLE || 'Discord voice bot'
  }
}
