        _, err = l.svcCtx.CoreRpc.CreateApi(l.ctx, &core.ApiInfo{
		Path:        pointy.GetPointer("/{{.modelNameSnake}}/create"),
		Description: pointy.GetPointer("apiDesc.create{{.modelName}}"),
		ApiGroup:    pointy.GetPointer("{{.modelNameSnake}}"),
		Method:      pointy.GetPointer("POST"),
	})

	if err != nil {
		return err
	}

	_, err = l.svcCtx.CoreRpc.CreateApi(l.ctx, &core.ApiInfo{
		Path:        pointy.GetPointer("/{{.modelNameSnake}}/update"),
		Description: pointy.GetPointer("apiDesc.update{{.modelName}}"),
		ApiGroup:    pointy.GetPointer("{{.modelNameSnake}}"),
		Method:      pointy.GetPointer("POST"),
	})

	if err != nil {
		return err
	}

	_, err = l.svcCtx.CoreRpc.CreateApi(l.ctx, &core.ApiInfo{
		Path:        pointy.GetPointer("/{{.modelNameSnake}}/delete"),
		Description: pointy.GetPointer("apiDesc.delete{{.modelName}}"),
		ApiGroup:    pointy.GetPointer("{{.modelNameSnake}}"),
		Method:      pointy.GetPointer("POST"),
	})

	if err != nil {
		return err
	}

	_, err = l.svcCtx.CoreRpc.CreateApi(l.ctx, &core.ApiInfo{
		Path:        pointy.GetPointer("/{{.modelNameSnake}}/list"),
		Description: pointy.GetPointer("apiDesc.get{{.modelName}}List"),
		ApiGroup:    pointy.GetPointer("{{.modelNameSnake}}"),
		Method:      pointy.GetPointer("POST"),
	})

	if err != nil {
		return err
	}

	_, err = l.svcCtx.CoreRpc.CreateApi(l.ctx, &core.ApiInfo{
		Path:        pointy.GetPointer("/{{.modelNameSnake}}"),
		Description: pointy.GetPointer("apiDesc.get{{.modelName}}ById"),
		ApiGroup:    pointy.GetPointer("{{.modelNameSnake}}"),
		Method:      pointy.GetPointer("POST"),
	})

	if err != nil {
		return err
	}