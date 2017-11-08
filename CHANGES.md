# Change log - CloudControl plugins for Packer

## v0.1.3-beta7

* Don't send admin password when deploying a server if no password was provided (DimensionDataResearch/packer-plugins-ddcloud#3).

## v0.1.3-beta6

* Implement deployment of servers from customer images that don't support Guest OS Customisation (DimensionDataResearch/packer-plugins-ddcloud#4).

## v0.1.3-beta5

* Upgrade to latest version of CloudControl client to enable working with images that have Guest OS Customisation disabled (DimensionDataResearch/packer-plugins-ddcloud#3).  
  Note that further work will be required to complete the implementation (see DimensionDataResearch/packer-plugins-ddcloud#4 for details).

## v0.1.3-beta3

* Fix WinRM port bug (partially in Packer).

## v0.1.3-beta2

* Fix incorrect population of WinRM password.

## v0.1.3-beta1

* Use `initial_admin_password`, if specified.

## v0.1.3-alpha1

### Changes

This is an early release of the CloudControl plugins for Packer intended to gather feedback on the overall design and behaviour.
It is not ready for production use.

* Implement customer image import post-processor
* Add initial documentation
