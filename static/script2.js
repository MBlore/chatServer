$(() => {
    var btnSignup = document.getElementById("signup");
    var form = document.getElementById("form");

    enableForm = function() {
        // Re-enable the form.
        var elements = form.elements;
        for (var i = 0, len = elements.length; i < len; ++i) {
            elements[i].disabled = false;
        }

        btnSignup.innerText = "Sign Up"
        btnSignup.disabled = false;
    }

    $("#form").submit((e) => {
        e.preventDefault();

        $("#validate").hide();
       
        // Validate form.
        $("#email").attr("required", "true");
        $("#username").attr("required", "true");
        $("#password").attr("required", "true");
        $("#password2").attr("required", "true");
        $("#displayname").attr("required", "true");

        var isValid = true;
        if (!$("#email")[0].checkValidity()) isValid = false;
        if (!$("#username")[0].checkValidity()) isValid = false;
        if (!$("#password")[0].checkValidity()) isValid = false;
        if (!$("#displayname")[0].checkValidity()) isValid = false;

        if (!isValid) {
            $("#validate").text("Please correct the errors on the form.");
            $("#validate").show();
            return;
        }

        if ($("#password").val() != $("#password2").val()) {
            $("#validate").text("The confirm password field must match the password field.");
            $("#validate").show();
            return;
        }
		
		var postData = $("#form").serialize();

        // Disable the form.
        var elements = form.elements;
        for (var i = 0, len = elements.length; i < len; ++i) {
            elements[i].disabled = true;
        }

        btnSignup.innerText = "Signing up..."
        btnSignup.disabled = true;
		
        // Data post.
        $.ajax({
            url: "/dosignup",
            type: "POST",
            data: postData,
            success: (data, ts, hr) => {
                if (hr.responseJSON && hr.responseJSON.hasOwnProperty("result") && hr.responseJSON.result == true) {
                    $("#success").show();
                    $("#fs").hide();
                } else {
                    $("#error").text("Woops, something went wrong. Please try again.");
                    $("#error").show();

                    enableForm();
                }
            },
            error: (hr, ts, err) => {
                
                $("#error").text("Woops, something went wrong. Please try again.");

                if (hr.responseJSON && hr.responseJSON.hasOwnProperty("errorCode")) {
                    switch(hr.responseJSON.errorCode) {
                        case -1: $("#error").text("Woops, something went wrong. Please try again."); break;
                        case -2: $("#error").text("The username is invalid. Please enter a different username."); break;
                        case -3: $("#error").text("The username already exists. Please enter a different username."); break;
                        case -4: $("#error").text("The password is invalid. Please enter a different password."); break;
                        case -5: $("#error").text("The email address is invalid. Please enter a different email address."); break;
                        case -6: $("#error").text("The display name is invalid. Please enter a different display name."); break;
                    }
                }
                
                $("#error").show();
                
                enableForm();
            }
        });

        return false;
    });

    btnSignup.addEventListener("click", (e) => {
        e.preventDefault();
                        
        $("#form").submit();
    });
});