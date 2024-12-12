document.getElementById('report').addEventListener("submit", function (event) {
    var token = localStorage.getItem('token');
    event.preventDefault();

    var selectedRadio = document.querySelector('input[name="format"]:checked');

    if (selectedRadio) {
        var url = '/admin/reports/' + selectedRadio.value;

        SendRequest(url, token, selectedRadio.value);
    } else {
        alert('Please select a report format.');
    }
});

function SendRequest(url, token, fileType) {
    fetch(url, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + token,
        },
    })
        .then(response => {
            if (!response.ok) {
                throw new Error('Error: ' + response.statusText);
            }
            return response.blob();
        })
        .then(blob => {
            var fileName = fileType === 'pdf' ? 'report.pdf' : 'report.xlsx';
            saveAs(blob, fileName);
        })
        .catch(error => {
            console.error('Error downloading file:', error);
            alert('Error downloading file: ' + error.message); // User feedback
        })
        .finally(() => {
            // Hide loading indicator
            document.getElementById('loading').style.display = 'none';
        });
}

document.getElementById('export').addEventListener("click", function () {
    var token = localStorage.getItem('token');

    // Проверяем наличие токена
    if (!token) {
        console.error('Токен не найден. Пожалуйста, войдите в систему.');
        return;
    }

    fetch('/admin/export', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + token, // Исправлено здесь
        },
    })
        .then(response => {
            if (!response.ok) {
                return response.json().then(err => { // Изменено для получения сообщения об ошибке
                    throw new Error(err.message || 'Ошибка при загрузке файла');
                });
            }
            return response.blob();
        })
        .then(blob => {
            saveAs(blob, 'export.json');
        })
        .catch(error => {
            console.error('Ошибка при скачивании файла:', error);
        });
});


document.getElementById('import').addEventListener('submit', function (event) {
    event.preventDefault();

    var formData = new FormData(event.target);

    const token = localStorage.getItem('token');

    if (!token) {
        alert("No token found. Please log in.");
        return;
    }
    fetch('/admin/import', {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${token}`
        },
        body: formData
    })
        .then(response => {
            if(response.ok){
                alert('Success.');
            }
            else{

                // Если код ответа не 200, отображаем сообщение в alert
                return response.json().then(data => {
                    const errorMessage = data.message || 'Unknown error occurred';
                    throw new Error(errorMessage);
                });

            }
        })
        .catch(error => {
            console.error('Ошибка при отправке запроса:', error);
            alert('Ошибка: ' + error.message);
        });


});

