package frecovery

import (
	"sync"

	"gitee.com/zengtao321/frdocker/config"
	"gitee.com/zengtao321/frdocker/frecovery/algo"
	"github.com/robfig/cron/v3"
)

func (app *FrecoveryApp) monitorMetric() *cron.Cron {
	app.Logger.Info("start metric monitoring...")
	c := cron.New()
	c.AddFunc(config.CONTAINER_METRIC_MONITOR_INTERVAL, app.monitorMetricScheduledTask)
	c.Start()
	return c
}

func (app *FrecoveryApp) monitorMetricScheduledTask() {
	groupedContainers := app.getGroupedContainers()
	for _, containers := range groupedContainers {
		go app.monitorMetricForContainerGroup(containers)
	}
}

func (app *FrecoveryApp) getGroupedContainers() [][]*Container {
	containers := [][]*Container{}
	for _, service := range app.Services {
		ctns := []*Container{}
		for _, ctn := range service.Containers {
			ctns = append(ctns, app.GetContainer(ctn))
		}
		containers = append(containers, ctns)
	}
	return containers
}

func (app *FrecoveryApp) monitorMetricForContainerGroup(containers []*Container) {
	n := len(containers)
	var wg sync.WaitGroup
	wg.Add(n)
	data := make([][]float64, n)
	for idx, container := range containers {
		app.Pool.Submit(app.monitorTaskFunc(data, idx, container, &wg))
	}
	wg.Wait()
	allEcc, thresh := algo.CalculateWithSample(data)
	for idx, ecc := range allEcc {
		container := containers[idx]
		container.Monitor.UpdateContainerEcc(ecc, thresh)
		if ecc > thresh {
			app.Logger.Errorf("[metric][%s][%s:%d] ecc: %.4f, thresh: %.4f", container.ServiceName, container.IP, container.Port, ecc, thresh)
		}
	}
}

func (app *FrecoveryApp) monitorTaskFunc(data [][]float64, idx int, container *Container, wg *sync.WaitGroup) AntsTaskWrapper {
	return func() {
		defer wg.Done()
		container.Monitor.UpdateContainerMetric(app.DockerCli)
		metric := container.Monitor.Metric
		data[idx] = []float64{metric.CPU, metric.Mem, metric.NetUp, metric.NetDn, metric.DiskR, metric.DiskW}
	}
}
