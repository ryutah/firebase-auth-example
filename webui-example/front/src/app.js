import "./main.css";

import config from "./firebase_config";

import * as firebase from "firebase";
import * as firebaseui from "./firebaseui/npm__ja";

firebase.initializeApp(config);

{
  const uiConfig = {
    callbacks: {
      signInSuccess: () => false
    },
    signInOptions: [firebase.auth.EmailAuthProvider.PROVIDER_ID],
    credentialHelper: firebaseui.auth.CredentialHelper.NONE
  };

  ////////////////////////////////////
  // DOMの取得とイベントリスナーの設定
  ////////////////////////////////////
  const ui = new firebaseui.auth.AuthUI(firebase.auth());

  const signOutButton = document.querySelector("#sign-out");
  signOutButton.addEventListener("click", e => {
    e.preventDefault();
    firebase
      .auth()
      .signOut()
      .then(() => {})
      .catch(alert);
  });

  const emailVerifyButton = document.querySelector("#email-verify");
  emailVerifyButton.addEventListener("click", e => {
    e.preventDefault();
    firebase
      .auth()
      .currentUser.sendEmailVerification({ url: "http://localhost:8080" })
      .then(() => {
        emailVerifyElements.forEach(dom =>
          dom.setAttribute("style", "display: none")
        );
        alert("メールを送信しました");
      })
      .catch(alert);
  });

  const signInElements = document.querySelectorAll(".on-signin");

  const emailVerifyElements = document.querySelectorAll(".on-email-verify");

  const sendToServer = document.querySelector("#send");
  sendToServer.addEventListener("click", e => {
    e.preventDefault();
    sendSampleGet();
  });
  ////////////////////////////////////

  ////////////////////////////////////
  // 認証関連のDOM操作
  ////////////////////////////////////
  const onSignIn = user => {
    signInElements.forEach(dom => dom.setAttribute("style", ""));
    if (!isVerifyUserEmail(user)) {
      emailVerifyElements.forEach(dom => dom.setAttribute("style", ""));
    }
  };

  const onSignOut = () => {
    const signInElements = document.querySelectorAll(".on-signin");
    signInElements.forEach(dom => dom.setAttribute("style", "display: none"));
    ui.start("#firebaseui-auth-container", uiConfig);
  };

  const isVerifyUserEmail = user => {
    const provider = user.providerData.filter(p => p.providerId === "password");
    if (provider.length === 0) {
      return true;
    }
    return user.emailVerified;
  };
  ////////////////////////////////////

  ////////////////////////////////////
  // APIリクエスト
  ////////////////////////////////////
  const sendSampleGet = () => {
    firebase
      .auth()
      .currentUser.getIdToken(true)
      .then(idToken =>
        fetch("http://localhost:8081/sample", {
          method: "GET",
          headers: { Authorization: idToken }
        })
      )
      .then(response => response.json())
      .then(console.log)
      .catch(console.error);
  };

  ////////////////////////////////////
  // 認証コールバック
  ////////////////////////////////////
  firebase.auth().onAuthStateChanged(user => {
    if (user) {
      onSignIn(user);
    } else {
      onSignOut(user);
    }
  });
}
