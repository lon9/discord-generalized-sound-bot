export default function ({ store, error, redirect, route }) {
  if (!store.state.authUser) {
    return redirect('/login', { path: route.path })
  }
}