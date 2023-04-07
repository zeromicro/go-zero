        _, err = l.svcCtx.CoreRpc.CreateApi(l.ctx, &core.ApiInfo{
		Path:        "/{{.modelNameSnake}}/create",
		Description: "apiDesc.create{{.modelName}}",
		ApiGroup:    "{{.modelNameSnake}}",
		Method:      "POST",
	})

	if err != nil {
		return err
	}

	_, err = l.svcCtx.CoreRpc.CreateApi(l.ctx, &core.ApiInfo{
		Path:        "/{{.modelNameSnake}}/update",
		Description: "apiDesc.update{{.modelName}}",
		ApiGroup:    "{{.modelNameSnake}}",
		Method:      "POST",
	})

	if err != nil {
		return err
	}

	_, err = l.svcCtx.CoreRpc.CreateApi(l.ctx, &core.ApiInfo{
		Path:        "/{{.modelNameSnake}}/delete",
		Description: "apiDesc.delete{{.modelName}}",
		ApiGroup:    "{{.modelNameSnake}}",
		Method:      "POST",
	})

	if err != nil {
		return err
	}

	_, err = l.svcCtx.CoreRpc.CreateApi(l.ctx, &core.ApiInfo{
		Path:        "/{{.modelNameSnake}}/list",
		Description: "apiDesc.get{{.modelName}}List",
		ApiGroup:    "{{.modelNameSnake}}",
		Method:      "POST",
	})

	if err != nil {
		return err
	}

	_, err = l.svcCtx.CoreRpc.CreateApi(l.ctx, &core.ApiInfo{
		Path:        "/{{.modelNameSnake}}",
		Description: "apiDesc.get{{.modelName}}ById",
		ApiGroup:    "{{.modelNameSnake}}",
		Method:      "Post",
	})

	if err != nil {
		return err
	}