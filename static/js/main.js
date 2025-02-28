// Custom JavaScript for the GTD application

// This event fires when HTMX is ready
document.addEventListener('DOMContentLoaded', function() {
    console.log('GTD App initialized');
    
    // Initialize quick capture modal functionality
    initQuickCapture();
});

// This event fires after an HTMX request completes
document.body.addEventListener('htmx:afterRequest', function(event) {
    // If quick capture was successful, reset the form and focus
    if (event.detail.target && event.detail.target.id === 'quick-capture-result') {
        const form = document.getElementById('quick-capture-form');
        if (form) {
            setTimeout(() => {
                form.reset();
                document.querySelector('#quick-capture-form input[name="title"]').focus();
            }, 100);
        }
    }
});

// Initializes quick capture functionality
function initQuickCapture() {
    // Handle keyboard shortcuts
    document.addEventListener('keydown', function(e) {
        // Alt+N to open quick capture modal from anywhere
        if (e.altKey && e.key === 'n') {
            e.preventDefault();
            const modal = document.getElementById('quick-capture-modal');
            if (modal) {
                modal.showModal();
                // Focus the title input
                setTimeout(() => {
                    document.querySelector('#quick-capture-form input[name="title"]').focus();
                }, 100);
            }
        }
        
        // Escape key to close the modal
        if (e.key === 'Escape') {
            const modal = document.getElementById('quick-capture-modal');
            if (modal && modal.open) {
                modal.close();
            }
        }
    });
    
    // When modal is opened, focus the title input
    const modal = document.getElementById('quick-capture-modal');
    if (modal) {
        modal.addEventListener('showModal', function() {
            document.querySelector('#quick-capture-form input[name="title"]').focus();
        });
        
        // Clear result when modal is closed
        modal.addEventListener('close', function() {
            const result = document.getElementById('quick-capture-result');
            if (result) {
                result.innerHTML = '';
            }
        });
    }
}