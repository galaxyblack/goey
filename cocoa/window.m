#include "_cgo_export.h"
#include "cocoa.h"
#import <Cocoa/Cocoa.h>

#define STYLE_MASK                                                             \
	( NSTitledWindowMask | NSClosableWindowMask | NSMiniaturizableWindowMask | \
	  NSResizableWindowMask )

@implementation NSWindow ( Goey )

- (BOOL)windowShouldClose:(id)sender {
	if ( windowShouldClose( sender ) != 0 ) {
		return YES;
	}
	return NO;
}

@end

@interface MyWindowDelegate : NSObject <NSWindowDelegate>
- (void)windowWillClose:(NSNotification*)aNotification;
- (void)windowDidResize:(NSNotification*)aNotification;
@end

@implementation MyWindowDelegate

- (void)windowWillClose:(NSNotification*)notification {
	NSWindow* window = [notification object];
	windowWillClose( window );
}

- (void)windowDidResize:(NSNotification*)notification {
	NSWindow* window = [notification object];
	windowDidResize( window );
}

@end

void* windowNew( char const* title, unsigned width, unsigned height ) {
	// Make sure that we have a delegate
	static MyWindowDelegate* delegate = 0;
	if ( !delegate ) {
		delegate = [[MyWindowDelegate alloc] init];
	}

	NSString* appName = [[NSString alloc] initWithUTF8String:title];

	NSWindow* window =
	    [[[NSWindow alloc] initWithContentRect:NSMakeRect( 0, 0, width, height )
	                                 styleMask:STYLE_MASK
	                                   backing:NSBackingStoreBuffered
	                                     defer:NO] autorelease];
	[window cascadeTopLeftFromPoint:NSMakePoint( 20, 20 )];
	[window setDelegate:delegate];
	[window setTitle:appName];
	[window makeKeyAndOrderFront:nil];
	return window;
}

void windowClose( void* handle ) {
	NSWindow* window = handle;
	[window close];
}

nssize_t windowContentSize( void* handle ) {
	NSUInteger style = [(NSWindow*)handle styleMask];
	NSRect frame = [(NSWindow*)handle frame];
	frame = [NSWindow contentRectForFrameRect:frame styleMask:style];

    nssize_t ret = { frame.size.width, frame.size.height };
	return ret;
}

void windowMakeFirstResponder( void* handle1, void* handle2 ) {
	NSWindow* w = (NSWindow*)handle1;
	NSControl* c = (NSControl*)handle2;

	[w makeFirstResponder:c];
}

void windowSetMinSize( void* handle, int width, int height ) {
	NSWindow* w = (NSWindow*)handle;

	// Adjust size from content to outer frame
	NSRect frame = NSMakeRect( 0, 0, width, height );
	frame = [NSWindow frameRectForContentRect:frame styleMask:[w styleMask]];
	[w setMinSize:NSMakeSize( NSWidth( frame ), NSHeight( frame ) )];
}
