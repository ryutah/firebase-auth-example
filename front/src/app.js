import "./main.css";

import config from "./firebase_config";

import * as firebase from "firebase";
import * as firebaseui from "./firebaseui/npm__ja";

firebase.initializeApp(config);

{
  const ui = new firebaseui.auth.AuthUI(firebase.auth());

  const signArea = document.querySelector("#sign-area");

  const showEmailVerifyComponents = () => {
    const verifyArea = document.querySelector("#email-verify-area");
    const verify = document.createElement("button");
    verify.addEventListener("click", e => {
      e.preventDefault();

      firebase
        .auth()
        .currentUser.sendEmailVerification({ url: "http://localhost:8080/" })
        .then(() => alert("メールを送ったぜ"))
        .catch(console.error);
    });
    verify.innerText = "アドレス確認メール送信";
    verifyArea.appendChild(verify);
  };

  const showLogout = () => {
    const logoutBtn = document.createElement("button");
    logoutBtn.id = "logout";
    logoutBtn.innerText = "Logout";
    logoutBtn.addEventListener("click", e => {
      e.preventDefault();
      firebase
        .auth()
        .signOut()
        .then(() => {})
        .catch(console.error);
    });

    signArea.appendChild(logoutBtn);
  };
  const hideLogout = () => {
    const logoutBtn = document.querySelector("#logout");
    if (logoutBtn) signArea.removeChild(logoutBtn);
  };

  const isVerifyUserEmail = user => {
    const provider = user.providerData.filter(p => p.providerId === "password");
    if (provider.length === 0) {
      return true;
    }
    return user.emailVerified;
  };

  const uiConfig = {
    callbacks: {
      signInSuccess: () => false
    },
    signInOptions: [firebase.auth.EmailAuthProvider.PROVIDER_ID],
    credentialHelper: firebaseui.auth.CredentialHelper.NONE
  };

  firebase.auth().onAuthStateChanged(user => {
    if (user) {
      showLogout();
      if (!isVerifyUserEmail(user)) {
        showEmailVerifyComponents();
      }
    } else {
      hideLogout();
      ui.start("#firebaseui-auth-container", uiConfig);
    }
  });
}
