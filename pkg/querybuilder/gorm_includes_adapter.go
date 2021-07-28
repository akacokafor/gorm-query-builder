package querybuilder

import "log"

func (g *GormAdapter) applyIncludes(instance OptionsInterface) error {
	if len(g.includesWhitelist) == 0 {
		log.Print(instance.GetIncludes())
		for _, val := range instance.GetIncludes() {
			relationshipName := g.normalizeIncludeName(val)
			g.addRelationship(relationshipName)
			g.db.Preload(relationshipName)
		}
		return nil
	}

	for _, suppliedInclude := range instance.GetIncludes() {
		for _, whiteListIncludeEntry := range g.includesWhitelist {
			if _k, ok := whiteListIncludeEntry.(string); ok {
				if _k == suppliedInclude {
					relationshipName := g.normalizeIncludeName(_k)
					g.addRelationship(relationshipName)
					g.db.Preload(relationshipName)
				}
			}

			//if op, ok := whiteListFilterEntry.(GormAllowedFilter); ok {
			//	for _, _k := range op.Keys() {
			//		if _k == suppliedFilterKey {
			//			if err :=  op.Execute(g.db,instance); err != nil {
			//				return err
			//			}
			//		}
			//	}
			//}
		}
	}

	return nil
}

