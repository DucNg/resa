var input = document.getElementById('input');
input.onkeyup = function () {
	var filter = input.value.toUpperCase();
	
	//var table = document.getElementById("table");
	//var tr = table.getElementsByTagName("tr");
	var tr = table.getElementsByClassName("info");

	for (i = 0; i < tr.length; i++) {
		
		var nom = tr[i].getElementsByClassName("nom")[0];
		var prenom = tr[i].getElementsByClassName("prenom")[0];
		var mail = tr[i].getElementsByClassName("mail")[0];
		var phone = tr[i].getElementsByClassName("phone")[0];

		var foundNom = false;
		var foundPrenom = false;
		var foundMail = false;
		var foundPhone = false;

		if (nom) {
			foundNom = nom.innerHTML.toUpperCase().indexOf(filter) > -1;
		}
		if (prenom) {
			foundPrenom = prenom.innerHTML.toUpperCase().indexOf(filter) > -1;
		}
		if (mail) {
			foundMail = mail.innerHTML.toUpperCase().indexOf(filter) > -1;
		}
		if (phone) {
			foundPhone = phone.innerHTML.toUpperCase().indexOf(filter) > -1;
		}

		if (foundNom || foundPrenom || foundMail || foundPhone) {
			tr[i].style.display = "";
		} else {
			tr[i].style.display = "none";
		}
	}
}