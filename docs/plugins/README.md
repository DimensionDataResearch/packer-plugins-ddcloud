# Plugins

The Packer plugins for CloudControl fall into one of 2 categories:

* Builders  
A builder creates a virtual machine image.
* Post-processors  
A post-processor transforms or otherwise processes a virtual machine image.

## Builders

* [Customer image](builders/customerimage.md)  
The customer image builder deploys a new server in CloudControl, runs any configured provisioners against that server, then clones it to create a new customer image.

## Post-processors

* [Customer image export](postprocessors/customerimage-export.md)  
The customer image export post-processor exports (and optionally downloads) a customer image generated by the customer image builder.
* [Customer image import](postprocessors/customerimage-import.md)  
The customer image import post-processor converts a local VMWare (`.vmx`) virtual machine into OVF (`.ovf`) format, uploads it to CloudControl, and then imports it as a customer image.
