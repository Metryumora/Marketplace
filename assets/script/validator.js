/**
 * Created by Metr_yumora on 02.03.2017.
 */

minUsernameLength = 4;

function checkUsername() {
    if (document.getElementById("username").value == "") {
        document.getElementById("usernameWarn").innerHTML = "Username cannot be blank!";
        return false;
    }
    re = /^\w+$/;
    if (!re.test(document.getElementById("username").value)) {
        document.getElementById("usernameWarn").innerHTML = "Username must contain only letters, numbers and underscores!";
        return false;
    }
    if (document.getElementById("username").value.length < minUsernameLength) {
        document.getElementById("usernameWarn").innerHTML = "Username must contain at least " + minUsernameLength + " characters!";
        return false;
    }
    else document.getElementById("usernameWarn").innerHTML = "";
}

minPasswordLength = 8;

function checkPassword() {
    if (document.getElementById("pwd1").value.length < minPasswordLength) {
        document.getElementById("passWarn").innerHTML = "Password must contain at least " + minPasswordLength + " characters!";
        return false;
    }
    re = /^\w+$/;
    if (!re.test(document.getElementById("pwd1").value)) {
        document.getElementById("passWarn").innerHTML = "Password should not contain any specific symbols!";
        return false;
    }
    re = /[0-9]/;
    if (!re.test(document.getElementById("pwd1").value)) {
        document.getElementById("passWarn").innerHTML = "Password must contain at least one number (0-9)!";
        return false;
    }
    re = /[a-z]/;
    if (!re.test(document.getElementById("pwd1").value)) {
        document.getElementById("passWarn").innerHTML = "Password must contain at least one lowercase letter (a-z)!";
        return false;
    }
    re = /[A-Z]/;
    if (!re.test(document.getElementById("pwd1").value)) {
        document.getElementById("passWarn").innerHTML = "Password must contain at least one uppercase letter (A-Z)!";
        return false;
    }

    document.getElementById("passWarn").innerHTML = "";
    return true;
}

function confirmPassword() {
    if (document.getElementById("pwd1").value != document.getElementById("pwd2").value) {
        document.getElementById("confWarn").innerHTML = "Passwords do not match!";
        return false;
    }
    else document.getElementById("confWarn").innerHTML = "";
}


