function passwordCheck() {
	var pass1 = document.getElementById("pass1").value;
	var pass2 = document.getElementById("pass2").value;
	if (pass1 != pass2) {
		document.getElementById("pass1").style.borderColor = "#E34234";
		document.getElementById("pass2").style.borderColor = "#E34234";
		return false;
	}
	return true;
}

function addErrorMsg() {
	var node = document.createElement("P");
	var textnode = document.createTextNode("Les mots de passe ne correspondent pas");
	node.appendChild(textnode);
	document.getElementById("passwordBlock2").appendChild(node);
}