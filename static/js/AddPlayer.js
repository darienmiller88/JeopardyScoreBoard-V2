function enableForm(form, duration) {
    setTimeout(() => {
        const inputs = form.querySelectorAll('input, select, button');
        inputs.forEach(input => {
            input.disabled = false;
        });
        
        const submitBtn = form.querySelector('button[type="submit"]');
        submitBtn.innerText = submitBtn.dataset.originalText || '+ Add Player to Location';
        submitBtn.style.opacity = '1';
        submitBtn.style.cursor = 'pointer';
    }, duration);
}

// Disables the form to prevent multiple submissions. Of course, routes will be rate limited as well.
function disableForm(form, duration) {
    const inputs = form.querySelectorAll('input, select, button');
    inputs.forEach(input => {
        input.disabled = true;
    });
    
    const submitBtn = form.querySelector('button[type="submit"]');
    const originalText = submitBtn.innerHTML;
    submitBtn.innerHTML = '⏳ Adding...';
    submitBtn.style.opacity = '0.6';
    submitBtn.style.cursor = 'not-allowed';
    
    setTimeout(() => {
        inputs.forEach(input => {
            input.disabled = false;
        });
        submitBtn.innerHTML = originalText;
        submitBtn.style.opacity = '1';
        submitBtn.style.cursor = 'pointer';
    }, duration);
}

// Function to show toast notification for the client.
function showToast(message, type = 'success', duration = 2000) {
    const toast = document.getElementById('toast');
    
    const icons = {
        error: '✕',
        success: '✓',
        info: 'i',
        warning: '!'
    }; 
    
    const icon = icons[type] || '✓';
    
    toast.innerHTML = `
        <div class="toast-icon">${icon}</div>
        <span>${message}</span>
        <button class="toast-close" onclick="hideToast()">×</button>
    `;
    
    toast.className = 'toast ' + type;
    toast.classList.add('show');
    
    setTimeout(() => {
        toast.classList.remove('show');
    }, duration);
}

function hideToast() {
    const toast = document.getElementById('toast');
    toast.classList.remove('show');
}

function populateEditModal(button) {
    const playerId = button.dataset.playerId;
    const names = button.dataset.playerName.split(" ");
    const firstName = names[0]
    const lastName = names[1]
    
    //populate the edit player name form with the name of the player
    document.getElementById('edit-player-id').value = playerId;
    document.getElementById('edit-first-name').value = firstName;
    document.getElementById('edit-last-name').value = lastName;
    
    // Update the HTMX button with correct player ID
    const saveBtn = document.getElementById('confirm-edit-btn');
    saveBtn.setAttribute('hx-put', `/players/${playerId}`);
    saveBtn.setAttribute('hx-target', `#${playerId}`)
}

function populateDeleteModal(button) {
    const playerId = button.dataset.playerId;
    const playerName = button.dataset.playerName;
    
    //Add the player name to the delete player modal
    document.getElementById('delete-player-name').textContent = playerName;
    
    const deleteBtn = document.getElementById('confirm-delete-btn');
    
    // Update the delete button with correct URL and target htmx attributes
    deleteBtn.setAttribute('hx-delete', `/players/${playerId}`);
    deleteBtn.setAttribute('hx-target', `#${playerId}`)
}

function savePlayerEdit() {
    // Optional: Add validation or additional logic here
}