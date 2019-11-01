#include <internal/types.h>
#include <internal.h>

#include <github.com/ianmcmahon/fmsynth/stm32f469/display.h>
// type decl
// var  decl
// func decl
// const decl
// type def
// var  def
// func def
// init
void github$0$com$ianmcmahon$fmsynth$stm32f469$display$init() {
	static bool called = false;
	if (called) {
		return;
	}
	called = true;
	fmt$init();
	internal$init();
}
