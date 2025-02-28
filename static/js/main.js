// Custom JavaScript for the GTD application

// This event fires when HTMX is ready
document.addEventListener('DOMContentLoaded', function() {
    console.log('GTD App initialized');
});

// This event fires after an HTMX request completes
document.body.addEventListener('htmx:afterRequest', function(event) {
    console.log('HTMX request completed', event);
});