import axios from 'axios'
import Cookie from 'js-cookie'
import jwtDecode from 'jwt-decode'

let instance = axios.create({
  baseURL: process.env.baseUrl
})

export const state = () => ({
  authUser: null,
  token: null
})

export const mutations = {
  SET_TOKEN: function(state, token){
    state.token = token
  },
  SET_USER: function(state, user){
    state.authUser = user
  }
}

export const actions = {
  nuxtServerInit({ commit }, { req }) {
    try {
      const jwtCookie = req.headers.cookie.split(";").find(c => c.trim().startsWith("token="));
      if (jwtCookie) {
          let token = jwtCookie.split('=')[1]
          let payload = jwtDecode(token)
          let date = Date.now() / 1000
          if (payload.exp > date) {
              commit('SET_USER', payload)
              commit('SET_TOKEN', token)
          }
      }
    } catch (error) {

    }
  },
  async login({ commit }, { username, password }) {
    try {
      const { data } = await instance.post(process.env.baseUrl + '/login', { username, password })
      let payload = jwtDecode(data.token)
      Cookie.set('token', data.token, { expires: 1 / 24 * 6 })  // Expire for 6h       
      commit('SET_TOKEN', data.token)
      commit('SET_USER', payload)
    } catch (error) {
      if (error.response && error.response.status === 401) {
        throw new Error('Bad credentials')
      }
      throw error
    }
  },

  async logout({ commit }) {
    Cookie.remove('token')
    commit('SET_USER', null)
    commit('SET_TOKEN', null)
  }
}