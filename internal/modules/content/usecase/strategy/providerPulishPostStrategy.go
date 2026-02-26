package strategy

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IStrategy"
)

func NewPublishPostStrategy(h1 *PostExtensionStrategy, h2 *PostInsightStrategy, h3 *MediaStrategy, h4 *SettingStrategy, h5 *PostInsightStrategy) []IStrategy.IPublishPostStrategy {
	return []IStrategy.IPublishPostStrategy{h1, h2, h3, h4, h5}
}

func NewPublishDeleteStrategy(h1 *PostExtensionStrategy, h2 *PostInsightStrategy, h3 *MediaStrategy, h4 *SettingStrategy, h5 *PostInsightStrategy , h6 *EditLogStrategy) []IStrategy.IPublishDeleteStrategy {
	return []IStrategy.IPublishDeleteStrategy{h1, h2, h3, h4, h5, h6}
}

