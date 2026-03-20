//go:build darwin && cgo

package app

/*
#cgo darwin CFLAGS: -x objective-c -fobjc-arc
#cgo darwin LDFLAGS: -framework Cocoa

#import <Cocoa/Cocoa.h>

static void hideDockIcon(void) {
	dispatch_async(dispatch_get_main_queue(), ^{
		[NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
	});
}

static void activateAppNow(void) {
	if ([NSThread isMainThread]) {
		[NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
		[NSApp activateIgnoringOtherApps:YES];
		// Switch back to accessory after a short delay
		dispatch_after(dispatch_time(DISPATCH_TIME_NOW, (int64_t)(0.5 * NSEC_PER_SEC)), dispatch_get_main_queue(), ^{
			[NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
		});
	} else {
		dispatch_sync(dispatch_get_main_queue(), ^{
			[NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
			[NSApp activateIgnoringOtherApps:YES];
			dispatch_after(dispatch_time(DISPATCH_TIME_NOW, (int64_t)(0.5 * NSEC_PER_SEC)), dispatch_get_main_queue(), ^{
				[NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
			});
		});
	}
}
*/
import "C"

func hideDock() {
	C.hideDockIcon()
}

func activateApp() {
	C.activateAppNow()
}
