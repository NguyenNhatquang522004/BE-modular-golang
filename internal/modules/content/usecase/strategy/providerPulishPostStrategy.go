package strategy

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/modules/content/domain/IStrategy"
)

func NewPublishPostStrategy(h1 *PostExtensionStrategy, h2 *PostInsightStrategy, h3 *MediaStrategy, h4 *SettingStrategy) []IStrategy.IPublishPostStrategy {
	return []IStrategy.IPublishPostStrategy{h1, h2, h3, h4}
}
