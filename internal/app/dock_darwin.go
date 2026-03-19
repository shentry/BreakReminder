//go:build darwin && cgo

package app

/*
#cgo darwin CFLAGS: -x objective-c -fobjc-arc
#cgo darwin LDFLAGS: -framework Cocoa

#import <Cocoa/Cocoa.h>

static void hideDockIcon(void) {
	// Schedule on main queue so it runs after Fyne initializes NSApp
	dispatch_async(dispatch_get_main_queue(), ^{
		[NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
	});
}
*/
import "C"

func hideDock() {
	C.hideDockIcon()
}
