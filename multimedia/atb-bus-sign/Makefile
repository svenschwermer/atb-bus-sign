#
# Copyright (C) 2018 Sven Schwermer
#
# This is free software, licensed under the GNU General Public License v2.
#

include $(TOPDIR)/rules.mk

PKG_NAME:=atb-bus-sign
PKG_VERSION:=0.1
PKG_RELEASE:=1

PKG_MAINTAINER:=Sven Schwermer <sven@svenschwermer.de>
PKG_LICENSE:=GPL-2.0

include $(INCLUDE_DIR)/package.mk

define Package/atb-bus-sign
  SECTION:=multimedia
  CATEGORY:=Multimedia
  DEPENDS:=+libpthread +boost +boost-system +libopenssl
  TITLE:=AtB Bus Sign
endef

define Package/atb-bus-sign/description
	Bus sign for Trondheim based AtB bus stops.
endef

define Build/Prepare
	mkdir -p $(PKG_BUILD_DIR)
	$(CP) ./src/* $(PKG_BUILD_DIR)/
endef

define Package/atb-bus-sign/install
	$(INSTALL_DIR) $(1)/usr/bin
	$(INSTALL_BIN) $(PKG_BUILD_DIR)/atb-bus-sign $(1)/usr/bin/

	$(INSTALL_DIR) $(1)/etc/init.d
	$(INSTALL_BIN) ./files/atb-bus-sign.init $(1)/etc/init.d/atb-bus-sign
endef

$(eval $(call BuildPackage,atb-bus-sign))
