const API_URL = "http://localhost:8080/users";

// Load users on page load
$(document).ready(function () {
    loadUsers();
});

// Create User
$("#createUser").click(function () {
    const user = {
        name: $("#name").val(),
        email: $("#email").val(),
        password: $("#password").val()
    };

    $.ajax({
        url: API_URL,
        type: "POST",
        contentType: "application/json",
        data: JSON.stringify(user),
        success: function () {
            alert("User created!");
            loadUsers();
        },
        error: function (err) {
            alert(err.responseJSON.error);
        }
    });
});

// Fetch Users
function loadUsers() {
    $.get(API_URL, function (data) {
        let rows = "";
        data.forEach(user => {
            rows += `
                <tr>
                    <td>${user.id}</td>
                    <td>${user.name}</td>
                    <td>${user.email}</td>
                    <td>
                        <button class="delete" onclick="deleteUser(${user.id})">Delete</button>
                    </td>
                </tr>`;
        });
        $("#usersTable").html(rows);
    });
}

// Delete User
function deleteUser(id) {
    $.ajax({
        url: API_URL + "/" + id,
        type: "DELETE",
        success: function () {
            loadUsers();
        }
    });
}
