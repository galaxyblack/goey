#include "_cgo_export.h"
#include "cocoa.h"
#import <Cocoa/Cocoa.h>
#include <assert.h>

@interface GNOPThread : NSThread
- (void)main;
@end

@implementation GNOPThread

- (void)main {
	// Do nothing.  This is a NOP thread.
	return;
}

@end

static void detachAThread() {
	trace( __func__ );

	// We need to make sure that Cocoa is running multithreaded.  Otherwise,
	// use of autopool from other threads will not work propertly.  The notes
	// for NSAutoreleasePool indicate that we need to detach a thread to
	// cause this transition.
	NSThread* thread = [[GNOPThread alloc] init];
	[thread start];
	[thread release];
}

static void initApplication() {
	trace( __func__ );

	NSString* quitString = [[NSString alloc] initWithUTF8String:"Quit "];
	NSString* qString = [[NSString alloc] initWithUTF8String:"q"];

	[NSApplication sharedApplication];
	//[NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];

	id menubar = [[NSMenu new] autorelease];
	id appMenuItem = [[NSMenuItem new] autorelease];
	[menubar addItem:appMenuItem];
	[NSApp setMainMenu:menubar];
	id appMenu = [[NSMenu new] autorelease];
	id appName = [[NSProcessInfo processInfo] processName];
	id quitTitle = [quitString stringByAppendingString:appName];
	id quitMenuItem = [[[NSMenuItem alloc] initWithTitle:quitTitle
	                                              action:@selector( terminate: )
	                                       keyEquivalent:qString] autorelease];
	[appMenu addItem:quitMenuItem];
	[appMenuItem setSubmenu:appMenu];
}

static NSAutoreleasePool* pool = 0;

void init() {
	trace( __func__ );

	assert( !pool );
	assert( !NSApp );
	assert( [NSThread isMainThread] );

	detachAThread();
	assert( [NSThread isMultiThreaded] );

	// This is a global release pool that we will keep around.  This will still
	// cause leaks, but until we identify where the autoreleasepool is required,
	// this will get us running.
	pool = [[NSAutoreleasePool alloc] init];
	assert( pool );
	initApplication();
	assert( NSApp && ![NSApp isRunning] );
}

void run() {
	trace( __func__ );

	assert( [NSThread isMultiThreaded] );
	assert( [NSThread isMainThread] );
	assert( NSApp && ![NSApp isRunning] );
	assert( pool );

	[NSApp activateIgnoringOtherApps:YES];
	[NSApp run];
}

@interface DoStop : NSObject
- (void)main;
@end

@implementation DoStop

- (void)main {
	trace( __func__ );

	assert( [NSThread isMultiThreaded] );
	assert( [NSThread isMainThread] );
	assert( NSApp && [NSApp isRunning] );

	[NSApp stop:nil];
}

@end

void stop() {
	trace( __func__ );

	assert( [NSThread isMultiThreaded] );
	assert( [NSThread isMainThread] );
	assert( NSApp && [NSApp isRunning] );

	id thunk = [[DoStop alloc] init];
	[thunk performSelectorOnMainThread:@selector( main )
	                        withObject:nil
	                     waitUntilDone:NO];
	[thunk release];
}

@interface DoThunk : NSObject
- (void)main;
@end

@implementation DoThunk

- (void)main {
	trace( __func__ );

	assert( [NSThread isMultiThreaded] );
	assert( [NSThread isMainThread] );
	assert( NSApp && [NSApp isRunning] );

	callbackDo();
}

@end

void do_thunk() {
	trace( __func__ );

	assert( [NSThread isMultiThreaded] );
	assert( ![NSThread isMainThread] );
	assert( NSApp );

	while ( ![NSApp isRunning] ) {
		[NSThread sleepForTimeInterval:0.001];
	}

	id thunk = [[DoThunk alloc] init];
	[thunk performSelectorOnMainThread:@selector( main )
	                        withObject:nil
	                     waitUntilDone:YES];
	[thunk release];
}

bool_t isMainThread( void ) {
	trace( __func__ );

	return [NSThread isMainThread];
}
