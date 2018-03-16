import "./main.css";

import config from "./firebase_config";

import * as firebase from "firebase";
import * as firebaseui from "./firebaseui/npm__ja";

firebase.initializeApp(config);

const uiConfig = {
  signInOptions: [firebase.auth.EmailAuthProvider.PROVIDER_ID],
  credentialHelper: firebaseui.auth.CredentialHelper.NONE
};

const ui = new firebaseui.auth.AuthUI(firebase.auth());
ui.start("#firebaseui-auth-container", uiConfig);
