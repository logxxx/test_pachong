package stub

import "learn/mywebcrawler/module"

type ModuleInternal interface {
	module.Module
	IncrCalledCount()
	IncrAcceptedCount()
	IncrCompletedCount()
	IncrHandlingNumber()
	DecrHandlingNumber()
	Clear()
}
