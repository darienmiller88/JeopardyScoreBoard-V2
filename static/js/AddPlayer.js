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