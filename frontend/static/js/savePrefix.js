document.getElementById('saveButton').addEventListener('click', function() {
    const input = document.getElementById('prefixInput').value.trim();
    if (input === '') {
        // Show danger alert if input is blank
        document.getElementById('dangerAlert').style.display = 'block';
        document.getElementById('dangerAlert').style.opacity = 1; // Fade in the alert
        document.getElementById('successAlert').style.display = 'none';
    } else {
        // Show success alert if input is valid
        document.getElementById('successAlert').style.display = 'block';
        document.getElementById('successAlert').style.opacity = 1; // Fade in the alert
        document.getElementById('dangerAlert').style.display = 'none';
        const baseurl = document.body.getAttribute('data-url');
        const guildID = document.body.getAttribute('data-guild-id');
        const url = `${baseurl}/dashboard/${guildID}/manage/update-prefix`;
        console.log(input);
        fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ prefix: input }), // Send the input data as JSON
        })
        .then(response => response.text())  // Handle the response as plain text
        .then(data => {
            // Handle the plain text response here
            console.log('Response:', data);
        })
        .catch((error) => {
            // Handle any errors
            console.error('Error:', error);
        });
    }

    // Hide the alert after 3 seconds
    setTimeout(function() {
        document.getElementById('dangerAlert').style.opacity = 0;
        setTimeout(function() {
            document.getElementById('dangerAlert').style.display = 'none';
        }, 300); // Wait for fade-out transition to complete

        document.getElementById('successAlert').style.opacity = 0;
        setTimeout(function() {
            document.getElementById('successAlert').style.display = 'none';
        }, 300); // Wait for fade-out transition to complete
    }, 3000); // 3 seconds
});