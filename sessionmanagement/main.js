import firebase from "firebase/app"
import "firebase/auth"

const config = {
  apiKey: process.env.API_KEY,
  authDomain: `${process.env.PROJECT_ID}.firebaseapp.com`,
  projectId: process.env.PROJECT_ID,
}

firebase.initializeApp(config)

const $email = document.querySelector("input[name=email]")
const $pass = document.querySelector("input[name=pass]")

document.querySelector("#signin").addEventListener("click", () => {
  firebase.auth().signInWithEmailAndPassword($email.value, $pass.value)
    .then(result => result.user.getIdToken())
    .then(idToken => {
      const body = new URLSearchParams()
      body.append("idToken", idToken)
      return fetch("/signin", {
        credentials: "include",
        method: "POST",
        body: body,
      })
    })
    .then(response => {
      if (response.status != 201) throw "failed to create session"
      return response.json()
    })
    .then(() => alert("Sigin in success"))
    .catch(e => {
      console.error(e)
      alert("Sign in failed")
    })
})

document.querySelector("#signout").addEventListener("click", () => {
  firebase
    .auth()
    .signOut()
    .then(() => {
      return fetch("/signout", {
        credentials: "include",
        method: "POST",
      })
    })
    .then(response => {
      if (response.status != 200) throw "Sign out failed"
      return response.text()
    })
    .then(alert)
    .catch(alert);
})

document.querySelector("#verify").addEventListener("click", () => {
  fetch("/verify", {
    credentials: "include",
    method: "GET",
  }).then(response => {
    if (response.status != 200) throw "invalid session"
    return response.text()
  }).then(alert)
    .catch(alert)
})
