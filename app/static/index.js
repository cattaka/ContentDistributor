window.addEventListener('load', function () {
    var signOutButton = document.getElementById('sign-out');
    if (signOutButton) {
        signOutButton.onclick = function () { firebase.auth().signOut(); };
    }

    // FirebaseUI config.
    var uiConfig = {
        signInSuccessUrl: '/',
        signInOptions: [
            // Leave the lines as is for the providers you want to offer your users.
            firebase.auth.GoogleAuthProvider.PROVIDER_ID,
            // firebase.auth.FacebookAuthProvider.PROVIDER_ID,
            // firebase.auth.TwitterAuthProvider.PROVIDER_ID,
            // firebase.auth.GithubAuthProvider.PROVIDER_ID,
            // firebase.auth.EmailAuthProvider.PROVIDER_ID,
            // firebase.auth.PhoneAuthProvider.PROVIDER_ID
        ],
        // Terms of service url.
        tosUrl: 'http://'
    };

    firebase.auth().onAuthStateChanged(function (user) {
        if (user) {
            // User is signed in.
            document.getElementById('account-details').textContent =
                'Signed in as ' + user.displayName + ' (' + user.email + ')';
            user.getIdToken().then(function (accessToken) {
                var f = document.getElementById('signIn-form');
                if (f.getAttribute("data-enabled") === 'true') {
                    document.getElementById('token').value = accessToken;
                    f.submit();
                }
            });
        } else {
            var ui = new firebaseui.auth.AuthUI(firebase.auth());
            ui.start('#firebaseui-auth-container', uiConfig);
            document.getElementById('account-details').textContent = '';
            var f = document.getElementById('signOut-form');
            if (f.getAttribute("data-enabled") === 'true') {
                f.submit();
            }
        }
    }, function (error) {
        console.log(error);
        alert('Unable to log in: ' + error)
    });
});
