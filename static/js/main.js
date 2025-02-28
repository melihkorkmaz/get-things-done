// Custom JavaScript for the GTD application

// This event fires when HTMX is ready
document.addEventListener('DOMContentLoaded', function() {
    console.log('GTD App initialized');
    
    // Initialize quick capture modal functionality
    initQuickCapture();
    
    // Initialize responsive sidebar functionality
    initSidebar();
    
    // Set active menu item based on current URL
    highlightCurrentMenuItem();
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

// Initialize mobile sidebar functionality
function initSidebar() {
    // Add responsive classes to sidebar
    const sidebar = document.querySelector('aside');
    if (sidebar) {
        sidebar.classList.add('sidebar');
        
        // Create toggle button for mobile
        const toggleButton = document.createElement('button');
        toggleButton.className = 'btn btn-ghost btn-circle lg:hidden';
        toggleButton.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" /></svg>';
        toggleButton.id = 'sidebar-toggle';
        
        // Add button to header
        const header = document.querySelector('header .flex');
        if (header) {
            header.insertBefore(toggleButton, header.firstChild);
        }
        
        // Create overlay for closing sidebar
        const overlay = document.createElement('div');
        overlay.className = 'sidebar-overlay';
        sidebar.parentNode.insertBefore(overlay, sidebar.nextSibling);
        
        // Toggle sidebar on button click
        toggleButton.addEventListener('click', function() {
            sidebar.classList.toggle('open');
        });
        
        // Close sidebar when clicking overlay
        overlay.addEventListener('click', function() {
            sidebar.classList.remove('open');
        });
        
        // Close sidebar when clicking menu item on mobile
        const menuItems = sidebar.querySelectorAll('a');
        menuItems.forEach(item => {
            item.addEventListener('click', function() {
                if (window.innerWidth < 768) {
                    sidebar.classList.remove('open');
                }
            });
        });
    }
}

// Highlight current menu item based on URL
function highlightCurrentMenuItem() {
    const currentPath = window.location.pathname;
    const currentSearch = window.location.search;
    
    // Find all menu links
    const menuLinks = document.querySelectorAll('aside .menu a');
    
    menuLinks.forEach(link => {
        // Remove any active class first
        link.classList.remove('active');
        
        const linkPath = link.getAttribute('href').split('?')[0];
        const linkSearch = link.getAttribute('href').includes('?') ? 
            link.getAttribute('href').split('?')[1] : '';
        
        // Exact match for path and search query
        if (linkPath === currentPath && linkSearch === currentSearch.substring(1)) {
            link.classList.add('active');
        }
        // Match just for path if it's not the root path
        else if (linkPath === currentPath && linkPath !== '/' && linkPath !== '/tasks') {
            link.classList.add('active');
        }
    });
}